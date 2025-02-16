import React, { useEffect, useState } from 'react';

function PostsPage() {
  const [posts, setPosts] = useState([]);
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [error, setError] = useState('');

  const token = localStorage.getItem('token');

  const fetchPosts = () => {
    fetch('http://localhost:8080/posts', {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })
      .then(async (res) => {
        if (!res.ok) {
          const errText = await res.text();
          throw new Error(errText);
        }
        return res.json();
      })
      .then((data) => setPosts(data))
      .catch((err) => setError(err.message));
  };

  useEffect(() => {
    if (!token) {
      setError('Нет токена. Авторизуйтесь.');
      return;
    }
    fetchPosts();
  }, []);

  const handleAddPost = async (e) => {
    e.preventDefault();
    try {
      const res = await fetch('http://localhost:8080/posts', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ 
          user_id: 1,  // Пример: указываем ID автора (тут можно подставлять текущего юзера)
          title,
          content 
        }),
      });

      if (!res.ok) {
        const errText = await res.text();
        throw new Error(errText);
      }

      // Добавление успешно, заново получим список постов
      fetchPosts();
      setTitle('');
      setContent('');
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <div>
      <h2>Посты</h2>
      {error && <p style={{ color: 'red' }}>{error}</p>}

      <ul>
        {posts.map((post) => (
          <li key={post.id}>
            <strong>{post.title}</strong><br/>
            {post.content}
          </li>
        ))}
      </ul>

      <h3>Добавить пост</h3>
      <form onSubmit={handleAddPost}>
        <div>
          <label>Заголовок:</label><br/>
          <input
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            required
          />
        </div>
        <div>
          <label>Содержимое:</label><br/>
          <textarea
            value={content}
            onChange={(e) => setContent(e.target.value)}
            required
          />
        </div>
        <button type="submit">Сохранить</button>
      </form>
    </div>
  );
}

export default PostsPage;
