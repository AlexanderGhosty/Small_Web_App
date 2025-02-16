package database

import (
    "database/sql"
    "fmt"

    _ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
    connStr := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, fmt.Errorf("ошибка открытия соединения с БД: %w", err)
    }

    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("ошибка проверки соединения с БД: %w", err)
    }

    fmt.Println("Соединение с БД установлено")
    return db, nil
}

func CreateTables(db *sql.DB) error {
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100) NOT NULL,
            email VARCHAR(100) NOT NULL UNIQUE,
            password VARCHAR(200) NOT NULL
        );
        CREATE TABLE IF NOT EXISTS posts (
            id SERIAL PRIMARY KEY,
            user_id INT NOT NULL,
            title VARCHAR(200) NOT NULL,
            content TEXT,
            FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
        );
        CREATE TABLE IF NOT EXISTS comments (
            id SERIAL PRIMARY KEY,
            post_id INT NOT NULL,
            author VARCHAR(100) NOT NULL,
            text TEXT NOT NULL,
            FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE
        );
    `)
    if err != nil {
        return fmt.Errorf("ошибка создания таблиц: %w", err)
    }
    return nil
}
