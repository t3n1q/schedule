package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Не удалось подключиться к БД:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("БД не отвечает:", err)
	}
	fmt.Println("Успешное подключение к БД!")

	mux := http.NewServeMux()

	// CRUD
	mux.HandleFunc("/api/schedules", handleSchedules)
	mux.HandleFunc("/api/teachers", handleTeachers)
	mux.HandleFunc("/api/groups", handleGroups)
	mux.HandleFunc("/api/subjects", handleSubjects)
	mux.HandleFunc("/api/classrooms", handleClassrooms)

	// Оборачиваем mux в CORS‑middleware, чтобы разрешить запросы с фронтенда
	corsedMux := corsMiddleware(mux)

	fmt.Println("Сервер запущен на порту 8080...")
	log.Fatal(http.ListenAndServe(":8080", corsedMux))
}

//CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
}



// POST/PUT (числовые ID)
type Schedule struct {
	ID          int    `json:"id"`
	DayOfWeek   string `json:"dayOfWeek"`
	Timeslot    int    `json:"timeslot"`
	TeacherID   int    `json:"teacherId"`
	GroupID     int    `json:"groupId"`
	SubjectID   int    `json:"subjectId"`
	ClassroomID int    `json:"classroomId"`
}

// GET
type ScheduleView struct {
	ID        int    `json:"id"`
	DayOfWeek string `json:"dayOfWeek"`
	Timeslot  int    `json:"timeslot"`

	TeacherID   int    `json:"teacherId"`
	TeacherName string `json:"teacherName"`

	GroupID   int    `json:"groupId"`
	GroupName string `json:"groupName"`

	SubjectID   int    `json:"subjectId"`
	SubjectName string `json:"subjectName"`

	ClassroomID   int    `json:"classroomId"`
	ClassroomName string `json:"classroomName"`
}

type Teacher struct {
	ID       int    `json:"id"`
	FullName string `json:"fullName"`
}

type Group struct {
	ID        int    `json:"id"`
	GroupName string `json:"groupName"`
}

type Subject struct {
	ID          int    `json:"id"`
	SubjectName string `json:"subjectName"`
}

type Classroom struct {
	ID       int    `json:"id"`
	RoomName string `json:"roomName"`
}

//CRUD
func handleSchedules(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "GET":
		query := `
			SELECT
				s.id,
				s.day_of_week,
				s.timeslot,

				t.id AS teacher_id,
				t.full_name AS teacher_name,

				g.id AS group_id,
				g.group_name AS group_name,

				sub.id AS subject_id,
				sub.subject_name AS subject_name,

				c.id AS classroom_id,
				c.room_name AS classroom_name

			FROM schedules s
			JOIN teachers t ON s.teacher_id = t.id
			JOIN groups g ON s.group_id = g.id
			JOIN subjects sub ON s.subject_id = sub.id
			JOIN classrooms c ON s.classroom_id = c.id
			ORDER BY s.id
		`
		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, "Ошибка выборки расписания", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var schedules []ScheduleView
		for rows.Next() {
			var sv ScheduleView
			err := rows.Scan(
				&sv.ID,
				&sv.DayOfWeek,
				&sv.Timeslot,
				&sv.TeacherID, &sv.TeacherName,
				&sv.GroupID, &sv.GroupName,
				&sv.SubjectID, &sv.SubjectName,
				&sv.ClassroomID, &sv.ClassroomName,
			)
			if err != nil {
				http.Error(w, "Ошибка чтения расписания", http.StatusInternalServerError)
				return
			}
			schedules = append(schedules, sv)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(schedules)

	case "POST":
		// Создание новой записи расписания
		var s Schedule
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			http.Error(w, "Некорректные данные при создании расписания", http.StatusBadRequest)
			return
		}

		query := `
			INSERT INTO schedules(day_of_week, timeslot, teacher_id, group_id, subject_id, classroom_id)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id
		`
		err := db.QueryRow(query,
			s.DayOfWeek, s.Timeslot,
			s.TeacherID, s.GroupID,
			s.SubjectID, s.ClassroomID,
		).Scan(&s.ID)
		if err != nil {
			http.Error(w, "Ошибка добавления записи расписания", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s)

	case "PUT":
		// Обновление
		var s Schedule
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			http.Error(w, "Некорректные данные при обновлении расписания", http.StatusBadRequest)
			return
		}

		query := `
			UPDATE schedules
			SET day_of_week=$1, timeslot=$2,
			    teacher_id=$3, group_id=$4,
			    subject_id=$5, classroom_id=$6
			WHERE id=$7
		`
		_, err := db.Exec(query,
			s.DayOfWeek, s.Timeslot,
			s.TeacherID, s.GroupID,
			s.SubjectID, s.ClassroomID,
			s.ID,
		)
		if err != nil {
			http.Error(w, "Ошибка обновления записи расписания", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	case "DELETE":
		// Удаление
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "Не передан параметр id", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Некорректный id", http.StatusBadRequest)
			return
		}
		_, err = db.Exec(`DELETE FROM schedules WHERE id=$1`, id)
		if err != nil {
			http.Error(w, "Ошибка удаления записи расписания", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}


func handleTeachers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rows, err := db.Query(`SELECT id, full_name FROM teachers`)
		if err != nil {
			http.Error(w, "Ошибка выборки преподавателей", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var teachers []Teacher
		for rows.Next() {
			var t Teacher
			if err := rows.Scan(&t.ID, &t.FullName); err != nil {
				http.Error(w, "Ошибка чтения преподавателей", http.StatusInternalServerError)
				return
			}
			teachers = append(teachers, t)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(teachers)

	case "POST":
		var t Teacher
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			http.Error(w, "Некорректные данные при создании преподавателя", http.StatusBadRequest)
			return
		}
		query := `INSERT INTO teachers(full_name) VALUES($1) RETURNING id`
		err := db.QueryRow(query, t.FullName).Scan(&t.ID)
		if err != nil {
			http.Error(w, "Ошибка добавления преподавателя", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(t)

	case "PUT":
		var t Teacher
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			http.Error(w, "Некорректные данные при обновлении преподавателя", http.StatusBadRequest)
			return
		}
		_, err := db.Exec(`UPDATE teachers SET full_name=$1 WHERE id=$2`, t.FullName, t.ID)
		if err != nil {
			http.Error(w, "Ошибка обновления преподавателя", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	case "DELETE":
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "Не передан параметр id преподавателя", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Некорректный id преподавателя", http.StatusBadRequest)
			return
		}
		_, err = db.Exec(`DELETE FROM teachers WHERE id=$1`, id)
		if err != nil {
			http.Error(w, "Ошибка удаления преподавателя", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}


func handleGroups(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rows, err := db.Query(`SELECT id, group_name FROM groups`)
		if err != nil {
			http.Error(w, "Ошибка выборки групп", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var groups []Group
		for rows.Next() {
			var g Group
			if err := rows.Scan(&g.ID, &g.GroupName); err != nil {
				http.Error(w, "Ошибка чтения групп", http.StatusInternalServerError)
				return
			}
			groups = append(groups, g)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(groups)

	case "POST":
		var g Group
		if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
			http.Error(w, "Некорректные данные при создании группы", http.StatusBadRequest)
			return
		}
		query := `INSERT INTO groups(group_name) VALUES($1) RETURNING id`
		err := db.QueryRow(query, g.GroupName).Scan(&g.ID)
		if err != nil {
			http.Error(w, "Ошибка добавления группы", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(g)

	case "PUT":
		var g Group
		if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
			http.Error(w, "Некорректные данные при обновлении группы", http.StatusBadRequest)
			return
		}
		_, err := db.Exec(`UPDATE groups SET group_name=$1 WHERE id=$2`, g.GroupName, g.ID)
		if err != nil {
			http.Error(w, "Ошибка обновления группы", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	case "DELETE":
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "Не передан параметр id группы", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Некорректный id группы", http.StatusBadRequest)
			return
		}
		_, err = db.Exec(`DELETE FROM groups WHERE id=$1`, id)
		if err != nil {
			http.Error(w, "Ошибка удаления группы", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}


func handleSubjects(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rows, err := db.Query(`SELECT id, subject_name FROM subjects`)
		if err != nil {
			http.Error(w, "Ошибка выборки предметов", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var subjects []Subject
		for rows.Next() {
			var s Subject
			if err := rows.Scan(&s.ID, &s.SubjectName); err != nil {
				http.Error(w, "Ошибка чтения предметов", http.StatusInternalServerError)
				return
			}
			subjects = append(subjects, s)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(subjects)

	case "POST":
		var s Subject
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			http.Error(w, "Некорректные данные при создании предмета", http.StatusBadRequest)
			return
		}
		query := `INSERT INTO subjects(subject_name) VALUES($1) RETURNING id`
		err := db.QueryRow(query, s.SubjectName).Scan(&s.ID)
		if err != nil {
			http.Error(w, "Ошибка добавления предмета", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s)

	case "PUT":
		var s Subject
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			http.Error(w, "Некорректные данные при обновлении предмета", http.StatusBadRequest)
			return
		}
		_, err := db.Exec(`UPDATE subjects SET subject_name=$1 WHERE id=$2`, s.SubjectName, s.ID)
		if err != nil {
			http.Error(w, "Ошибка обновления предмета", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	case "DELETE":
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "Не передан параметр id предмета", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Некорректный id предмета", http.StatusBadRequest)
			return
		}
		_, err = db.Exec(`DELETE FROM subjects WHERE id=$1`, id)
		if err != nil {
			http.Error(w, "Ошибка удаления предмета", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}


func handleClassrooms(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rows, err := db.Query(`SELECT id, room_name FROM classrooms`)
		if err != nil {
			http.Error(w, "Ошибка выборки аудиторий", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var classrooms []Classroom
		for rows.Next() {
			var c Classroom
			if err := rows.Scan(&c.ID, &c.RoomName); err != nil {
				http.Error(w, "Ошибка чтения аудиторий", http.StatusInternalServerError)
				return
			}
			classrooms = append(classrooms, c)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(classrooms)

	case "POST":
		var c Classroom
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			http.Error(w, "Некорректные данные при создании аудитории", http.StatusBadRequest)
			return
		}
		query := `INSERT INTO classrooms(room_name) VALUES($1) RETURNING id`
		err := db.QueryRow(query, c.RoomName).Scan(&c.ID)
		if err != nil {
			http.Error(w, "Ошибка добавления аудитории", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(c)

	case "PUT":
		var c Classroom
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			http.Error(w, "Некорректные данные при обновлении аудитории", http.StatusBadRequest)
			return
		}
		_, err := db.Exec(`UPDATE classrooms SET room_name=$1 WHERE id=$2`, c.RoomName, c.ID)
		if err != nil {
			http.Error(w, "Ошибка обновления аудитории", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	case "DELETE":
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "Не передан параметр id аудитории", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Некорректный id аудитории", http.StatusBadRequest)
			return
		}
		_, err = db.Exec(`DELETE FROM classrooms WHERE id=$1`, id)
		if err != nil {
			http.Error(w, "Ошибка удаления аудитории", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
