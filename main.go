package main

import (
	"database/sql"
    "fmt"
    "log"
    "net/http"
    "encoding/json"

	_ "github.com/lib/pq"
)

// Var to store DB connection
var db *sql.DB

type User struct{
	ID int `json: "id"`
	Name string `json:"name"`
	Email string `json:"email"`
}

type Post struct {
	ID int `json:"id"`
	UserID int `json:"user_id"`
	Title string `json: "title"`
	Content string `json: "content"`
}

type Comment struct {
    ID     int    `json:"id"`
    PostID int    `json:"post_id"`
    Author string `json:"author"`
    Text   string `json:"text"`
}

func initDB(){
	connStr :="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Ошибка открытия соединения с БД: %v\n", err)
	}

	// Connection check 
	err = db.Ping()
	if err != nil {
		log.Fatalf("Ошибка проверки соединения с БД: %v\n", err)
	}
	
	fmt.Println("Соединение с БД установлено")

    // Создадим таблицы, если они не существуют
    createTables()
}

func createTables(){
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(100) NOT NULL UNIQUE
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
        log.Printf("Ошибка создания таблиц: %v\n", err)
    }
}


// --- USER handlers
func getUsersHandler(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT id, name, email FROM users")
    if err != nil {
        http.Error(w, "Ошибка при запросе пользователей", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var users []User
    for rows.Next() {
        var u User
        if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
            http.Error(w, "Ошибка чтения данных пользователя", http.StatusInternalServerError)
            return
        }
        users = append(users, u)
    }

    // Возвращаем в формате JSON
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
    var newUser User
    if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
        http.Error(w, "Некорректные данные пользователя", http.StatusBadRequest)
        return
    }

    query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`
    err := db.QueryRow(query, newUser.Name, newUser.Email).Scan(&newUser.ID)
    if err != nil {
        http.Error(w, "Ошибка при создании пользователя", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(newUser)
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
    var u User
    if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
        http.Error(w, "Некорректные данные для обновления", http.StatusBadRequest)
        return
    }

    query := `UPDATE users SET name = $1, email = $2 WHERE id = $3`
    _, err := db.Exec(query, u.Name, u.Email, u.ID)
    if err != nil {
        http.Error(w, "Ошибка при обновлении пользователя", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(u)
}

