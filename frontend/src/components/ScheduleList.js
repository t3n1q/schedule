import React, { useEffect, useState } from 'react';
import { getSchedules, deleteSchedule } from '../api';

const ScheduleList = ({ mode }) => {
  const [schedules, setSchedules] = useState([]); 

  useEffect(() => {
    fetchSchedules();
  }, []);

  const fetchSchedules = async () => {
    try {
      const data = await getSchedules();
      setSchedules(data); 
    } catch (error) {
      console.error('Ошибка при загрузке расписания:', error);
      setSchedules([]);
    }
  };

  const handleDelete = async (id) => {
    try {
      await deleteSchedule(id);
   
      setSchedules((prev) => prev.filter((item) => item.id !== id));
    } catch (error) {
      console.error('Ошибка при удалении расписания:', error);
    }
  };


  return (
    <table className="table table-bordered">
      <thead>
        <tr>
          <th>ID</th>
          <th>День</th>
          <th>Пара</th>
          <th>Преподаватель</th>
          <th>Группа</th>
          <th>Предмет</th>
          <th>Аудитория</th>
          <th>Действия</th>
        </tr>
      </thead>
      <tbody>
        {schedules.map((schedule) => (
          <tr key={schedule.id}>
            <td>{schedule.id}</td>
            <td>{schedule.dayOfWeek}</td>
            <td>{schedule.timeslot}</td>
            <td>{schedule.teacherName}</td>
            <td>{schedule.groupName}</td>
            <td>{schedule.subjectName}</td>
            <td>{schedule.classroomName}</td>
            <td>
              <button
                className="btn btn-danger btn-sm"
                onClick={() => handleDelete(schedule.id)}
              >
                Удалить
              </button>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
};

export default ScheduleList;
