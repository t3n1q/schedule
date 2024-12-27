import React from 'react';

const NavigationBar = ({ onSwitchView }) => {
  return (
    <div className="btn-group mb-3">
      <button
        className="btn btn-primary"
        onClick={() => onSwitchView('groups')}
      >
        По группам
      </button>
      <button
        className="btn btn-secondary"
        onClick={() => onSwitchView('teachers')}
      >
        По преподавателям
      </button>
    </div>
  );
};

export default NavigationBar;
