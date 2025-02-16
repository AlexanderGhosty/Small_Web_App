import React from 'react';
import { Link } from 'react-router-dom';

function HomePage() {
  return (
    <div style={{ maxWidth: '600px', margin: '0 auto' }}>
      <h1>Добро пожаловать в Small Web App</h1>
      <p>
        Это небольшое приложение с бэкендом на Go и фронтендом на React. 
        Ниже представлены основные разделы:
      </p>
      <ul>
        <li><Link to="/login">Авторизация (Login)</Link></li>
        <li><Link to="/users">Пользователи (Users)</Link></li>
        <li><Link to="/posts">Посты (Posts)</Link></li>
        <li><Link to="/comments">Комментарии (Comments)</Link></li>
      </ul>
    </div>
  );
}

export default HomePage;
