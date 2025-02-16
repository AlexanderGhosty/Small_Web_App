import React, { useEffect, useState } from 'react';

function CommentsPage() {
  const [comments, setComments] = useState([]);
  const [postID, setPostID] = useState('');
  const [author, setAuthor] = useState('');
  const [text, setText] = useState('');
  const [error, setError] = useState('');

  // Токен из localStorage (выданный бэкендом при логине)
  const token = localStorage.getItem('token');

  // Функция для получения всех комментариев
  const fetchComments = async () => {
    if (!token) {
      setError('Не найден токен. Авторизуйтесь, пожалуйста.');
      return;
    }

    try {
      const res = await fetch('http://localhost:8080/comments', {
        method: 'GET',
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!res.ok) {
        const errText = await res.text();
        throw new Error(errText);
      }

      const data = await res.json();
      setComments(data);
    } catch (err) {
      setError(err.message);
    }
  };

  // При первом рендере вызываем fetchComments
  useEffect(() => {
    fetchComments();
  }, []);

  // Функция для добавления комментария
  const handleAddComment = async (e) => {
    e.preventDefault();
    setError('');

    if (!postID || !author || !text) {
      setError('Все поля обязательны для заполнения');
      return;
    }

    try {
      const res = await fetch('http://localhost:8080/comments', {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          post_id: Number(postID),
          author,
          text,
        }),
      });

      if (!res.ok) {
        const errText = await res.text();
        throw new Error(errText);
      }

      // Добавление прошло успешно, обновляем список комментариев
      await fetchComments();

      // Очищаем поля формы
      setPostID('');
      setAuthor('');
      setText('');
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <div style={{ maxWidth: '600px', margin: '0 auto' }}>
      <h2>Комментарии</h2>

      {error && <p style={{ color: 'red' }}>{error}</p>}

      <ul>
        {comments.map((comment) => (
          <li key={comment.id} style={{ marginBottom: '1rem' }}>
            <strong>Автор:</strong> {comment.author} <br />
            <strong>Текст:</strong> {comment.text} <br />
            <strong>PostID:</strong> {comment.post_id}
          </li>
        ))}
      </ul>

      <h3>Добавить комментарий</h3>
      <form onSubmit={handleAddComment}>
        <div>
          <label>PostID:</label><br />
          <input
            type="number"
            value={postID}
            onChange={(e) => setPostID(e.target.value)}
            required
          />
        </div>
        <div>
          <label>Автор:</label><br />
          <input
            type="text"
            value={author}
            onChange={(e) => setAuthor(e.target.value)}
            required
          />
        </div>
        <div>
          <label>Текст комментария:</label><br />
          <textarea
            value={text}
            onChange={(e) => setText(e.target.value)}
            required
          />
        </div>
        <button type="submit" style={{ marginTop: '10px' }}>
          Добавить
        </button>
      </form>
    </div>
  );
}

export default CommentsPage;
