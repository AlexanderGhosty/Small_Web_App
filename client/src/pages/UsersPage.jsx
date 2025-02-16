import React, { useEffect, useState } from 'react';

function UsersPage() {
  const [users, setUsers] = useState([]);
  const [error, setError] = useState('');

  useEffect(() => {
    const token = localStorage.getItem('token');
    if (!token) {
      // Если токена нет — возможно, редирект на /login
      setError('Требуется авторизация');
      return;
    }

    fetch('http://localhost:8080/users', {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
      .then(async (res) => {
        if (!res.ok) {
          const errText = await res.text();
          throw new Error(errText);
        }
        return res.json();
      })
      .then((data) => setUsers(data))
      .catch((err) => setError(err.message));
  }, []);

  if (error) {
    return <p style={{color: 'red'}}>Ошибка: {error}</p>;
  }

  return (
    <div>
      <h2>Пользователи</h2>
      <ul>
        {users.map((user) => (
          <li key={user.id}>
            {user.name} ({user.email}) 
          </li>
        ))}
      </ul>
    </div>
  );
}

export default UsersPage;
