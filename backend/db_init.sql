CREATE TABLE IF NOT EXISTS teachers (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS groups (
    id SERIAL PRIMARY KEY,
    group_name VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS subjects (
    id SERIAL PRIMARY KEY,
    subject_name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS classrooms (
    id SERIAL PRIMARY KEY,
    room_name VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS schedules (
    id SERIAL PRIMARY KEY,
    day_of_week VARCHAR(20) NOT NULL,
    timeslot INTEGER NOT NULL,
    teacher_id INT NOT NULL,
    group_id INT NOT NULL,
    subject_id INT NOT NULL,
    classroom_id INT NOT NULL,
    CONSTRAINT fk_teacher
      FOREIGN KEY(teacher_id) REFERENCES teachers(id) ON DELETE CASCADE,
    CONSTRAINT fk_group
      FOREIGN KEY(group_id) REFERENCES groups(id) ON DELETE CASCADE,
    CONSTRAINT fk_subject
      FOREIGN KEY(subject_id) REFERENCES subjects(id) ON DELETE CASCADE,
    CONSTRAINT fk_classroom
      FOREIGN KEY(classroom_id) REFERENCES classrooms(id) ON DELETE CASCADE
);

-- Пример тестовых данных
INSERT INTO teachers (full_name) VALUES ('Куликов И.А.'), ('Сидоров П.П.');
INSERT INTO groups (group_name) VALUES ('Группа А-1'), ('Группа B-1');
INSERT INTO subjects (subject_name) VALUES ('Математика'), ('Физика');
INSERT INTO classrooms (room_name) VALUES ('Ауд. 101'), ('Ауд. 102');

INSERT INTO schedules(day_of_week, timeslot, teacher_id, group_id, subject_id, classroom_id)
VALUES
('Понедельник', 1, 1, 1, 1, 1),
('Понедельник', 2, 1, 1, 2, 2),
('Вторник',     1, 2, 2, 1, 1),
('Вторник',     2, 2, 2, 2, 2);
