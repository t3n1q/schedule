import React, { useState } from 'react';
import { createSchedule } from '../api';

const ScheduleForm = () => {
  const [dayOfWeek, setDayOfWeek] = useState('');
  const [timeslot, setTimeslot] = useState('');
  const [teacherId, setTeacherId] = useState('');
  const [groupId, setGroupId] = useState('');
  const [subjectId, setSubjectId] = useState('');
  const [classroomId, setClassroomId] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await createSchedule({
        dayOfWeek,
        timeslot: parseInt(timeslot),
        teacherId: parseInt(teacherId),
        groupId: parseInt(groupId),
        subjectId: parseInt(subjectId),
        classroomId: parseInt(classroomId),
      });
      alert('Запись успешно добавлена!');
      // Сброс
      setDayOfWeek('');
      setTimeslot('');
      setTeacherId('');
      setGroupId('');
      setSubjectId('');
      setClassroomId('');
    } catch (error) {
      console.error('Ошибка при создании записи:', error);
      alert('Ошибка при создании записи');
    }
  };

  return (
    <form onSubmit={handleSubmit} className="row g-3">
      <div className="col-md-2">
        <label className="form-label">День недели</label>
        <input
          type="text"
          className="form-control"
          value={dayOfWeek}
          onChange={(e) => setDayOfWeek(e.target.value)}
          required
        />
      </div>
      <div className="col-md-1">
        <label className="form-label">№ пары</label>
        <input
          type="number"
          className="form-control"
          value={timeslot}
          onChange={(e) => setTimeslot(e.target.value)}
          required
        />
      </div>
      <div className="col-md-2">
        <label className="form-label">Преподаватель (ID)</label>
        <input
          type="number"
          className="form-control"
          value={teacherId}
          onChange={(e) => setTeacherId(e.target.value)}
          required
        />
      </div>
      <div className="col-md-2">
        <label className="form-label">Группа (ID)</label>
        <input
          type="number"
          className="form-control"
          value={groupId}
          onChange={(e) => setGroupId(e.target.value)}
          required
        />
      </div>
      <div className="col-md-2">
        <label className="form-label">Предмет (ID)</label>
        <input
          type="number"
          className="form-control"
          value={subjectId}
          onChange={(e) => setSubjectId(e.target.value)}
          required
        />
      </div>
      <div className="col-md-2">
        <label className="form-label">Аудитория (ID)</label>
        <input
          type="number"
          className="form-control"
          value={classroomId}
          onChange={(e) => setClassroomId(e.target.value)}
          required
        />
      </div>
      <div className="col-md-1 d-flex align-items-end">
        <button type="submit" className="btn btn-success">
          Добавить
        </button>
      </div>
    </form>
  );
};

export default ScheduleForm;
