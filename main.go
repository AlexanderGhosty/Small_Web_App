package main

import (
	"database/sql"
    "fmt"
    "log"
    "net/http"
    "encoding/json"
    "time"
    "context"

	// Auntification
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

var jwtKey = []byte("SUPER_SECRET_KEY")

type Claims struct {
    UserID int    `json:"user_id"`
    Email  string `json:"email"`
    jwt.RegisteredClaims
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
        return
    }

    var creds struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
        http.Error(w, "Некорректные данные", http.StatusBadRequest)
        return
    }

    var user User
    var hashedPassword string
    err := db.QueryRow("SELECT id, name, email, password FROM users WHERE email = $1", creds.Email).
        Scan(&user.ID, &user.Name, &user.Email, &hashedPassword)
    if err != nil {
        http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password)); err != nil {
        http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
        return
    }

    expirationTime := time.Now().Add(60 * time.Minute)
    claims := &Claims{
        UserID: user.ID,
        Email:  user.Email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        http.Error(w, "Ошибка при генерации токена", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}



// Var to store DB connection
var db *sql.DB

type User struct{
	ID int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
}

type Post struct {
	ID int `json:"id"`
	UserID int `json:"user_id"`
	Title string `json:"title"`
	Content string `json:"content"`
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

    // Return in JSON format
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
		log.Printf("Ошибка декодирования JSON: %v\n", err)
        http.Error(w, "Некорректные данные для обновления", http.StatusBadRequest)
        return
    }

    query := `UPDATE users SET name = $1, email = $2 WHERE id = $3`
    _, err := db.Exec(query, u.Name, u.Email, u.ID)
    if err != nil {
		log.Printf("Ошибка обновления пользователя: %v\n", err)
        http.Error(w, "Ошибка при обновлении пользователя", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(u)
}

// ========== Handlers for Post (example) ==========

// Get all posts (GET)
func getPostsHandler(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT id, user_id, title, content FROM posts")
    if err != nil {
        http.Error(w, "Error retrieving posts", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var posts []Post
    for rows.Next() {
        var p Post
        if err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content); err != nil {
            http.Error(w, "Error reading post data", http.StatusInternalServerError)
            return
        }
        posts = append(posts, p)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(posts)
}

// Create a new post (POST)
func createPostHandler(w http.ResponseWriter, r *http.Request) {
    var p Post
    if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
        http.Error(w, "Invalid post data", http.StatusBadRequest)
        return
    }

    query := `INSERT INTO posts (user_id, title, content) VALUES ($1, $2, $3) RETURNING id`
    err := db.QueryRow(query, p.UserID, p.Title, p.Content).Scan(&p.ID)
    if err != nil {
        http.Error(w, "Error creating post", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(p)
}

// Update a post (PUT)
func updatePostHandler(w http.ResponseWriter, r *http.Request) {
    var p Post
    if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
        http.Error(w, "Invalid data for updating post", http.StatusBadRequest)
        return
    }

    query := `UPDATE posts SET title = $1, content = $2 WHERE id = $3`
    _, err := db.Exec(query, p.Title, p.Content, p.ID)
    if err != nil {
        http.Error(w, "Error updating post", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(p)
}

// ========== Handlers for Comment (example) ==========

// Get all comments (GET)
func getCommentsHandler(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT id, post_id, author, text FROM comments")
    if err != nil {
        http.Error(w, "Error retrieving comments", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var comments []Comment
    for rows.Next() {
        var c Comment
        if err := rows.Scan(&c.ID, &c.PostID, &c.Author, &c.Text); err != nil {
            http.Error(w, "Error reading comment data", http.StatusInternalServerError)
            return
        }
        comments = append(comments, c)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(comments)
}

// Create a new comment (POST)
func createCommentHandler(w http.ResponseWriter, r *http.Request) {
    var c Comment
    if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
        http.Error(w, "Invalid comment data", http.StatusBadRequest)
        return
    }

    query := `INSERT INTO comments (post_id, author, text) VALUES ($1, $2, $3) RETURNING id`
    err := db.QueryRow(query, c.PostID, c.Author, c.Text).Scan(&c.ID)
    if err != nil {
        http.Error(w, "Error creating comment", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(c)
}

// Update a comment (PUT)
func updateCommentHandler(w http.ResponseWriter, r *http.Request) {
    var c Comment
    if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
        http.Error(w, "Invalid data for updating comment", http.StatusBadRequest)
        return
    }

    query := `UPDATE comments SET author = $1, text = $2 WHERE id = $3`
    _, err := db.Exec(query, c.Author, c.Text, c.ID)
    if err != nil {
        http.Error(w, "Error updating comment", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(c)
}

func main (){
	// DB init
	initDB()

	// Rout registraton
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            getUsersHandler(w, r)        // Get all users
        case http.MethodPost:
            createUserHandler(w, r)      // Creat user
        case http.MethodPut:
            updateUserHandler(w, r)      // Update user
        default:
            http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
        }
    })

	http.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            getPostsHandler(w, r)        // Get all posts
        case http.MethodPost:
            createPostHandler(w, r)      // Create post
        case http.MethodPut:
            updatePostHandler(w, r)      // Update post
        default:
            http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
        }
    })

    http.HandleFunc("/comments", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            getCommentsHandler(w, r)     // Get all comments
        case http.MethodPost:
            createCommentHandler(w, r)   // Create comment
        case http.MethodPut:
            updateCommentHandler(w, r)   // Update comment
        default:
            http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
        }
    })
	
	fmt.Println("Сервер запущен на http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}