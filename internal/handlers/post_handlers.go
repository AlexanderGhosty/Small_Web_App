package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"

    "small_web_app/internal/models"
)

// Маршрутизатор для /posts
func PostHandlers(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        getPostsHandler(db, w, r)
    case http.MethodPost:
        createPostHandler(db, w, r)
    case http.MethodPut:
        updatePostHandler(db, w, r)
    default:
        http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
    }
}

func getPostsHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT id, user_id, title, content FROM posts")
    if err != nil {
        http.Error(w, "Error retrieving posts", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var posts []models.Post
    for rows.Next() {
        var p models.Post
        if err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content); err != nil {
            http.Error(w, "Error reading post data", http.StatusInternalServerError)
            return
        }
        posts = append(posts, p)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(posts)
}

func createPostHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    var p models.Post
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

func updatePostHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    var p models.Post
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
