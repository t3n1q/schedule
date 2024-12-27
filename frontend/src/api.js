const API_URL = 'http://localhost:8080/api';

export async function getSchedules() {
  const response = await fetch(`${API_URL}/schedules`);
  if (!response.ok) {
    throw new Error('Ошибка при получении расписания');
  }
  return response.json();
}

export async function createSchedule(schedule) {
  const response = await fetch(`${API_URL}/schedules`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(schedule),
  });
  if (!response.ok) {
    throw new Error('Ошибка при создании записи');
  }
  return response.json();
}

export async function deleteSchedule(id) {
  const response = await fetch(`${API_URL}/schedules?id=${id}`, {
    method: 'DELETE',
  });
  if (!response.ok) {
    throw new Error('Ошибка при удалении записи');
  }
  return response;
}
