import React, { useState } from 'react';
import NavigationBar from './components/NavigationBar';
import ScheduleList from './components/ScheduleList';
import ScheduleForm from './components/ScheduleForm';

function App() {
  const [view, setView] = useState('groups'); 

  const handleSwitchView = (newView) => {
    setView(newView);
  };

  return (
    <div className="container my-3">
      <h1>Система расписания</h1>
      <NavigationBar onSwitchView={handleSwitchView} />

      {view === 'groups' && (
        <>
          <h2>Расписание по группам</h2>
          <ScheduleList mode="groups" />
        </>
      )}
      {view === 'teachers' && (
        <>
          <h2>Расписание по преподавателям</h2>
          <ScheduleList mode="teachers" />
        </>
      )}

      <hr />
      <ScheduleForm />
    </div>
  );
}

export default App;
