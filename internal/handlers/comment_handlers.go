package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"

    "small_web_app/internal/models"
)

func CommentHandlers(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        getCommentsHandler(db, w, r)
    case http.MethodPost:
        createCommentHandler(db, w, r)
    case http.MethodPut:
        updateCommentHandler(db, w, r)
    default:
        http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
    }
}

func getCommentsHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT id, post_id, author, text FROM comments")
    if err != nil {
        http.Error(w, "Error retrieving comments", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var comments []models.Comment
    for rows.Next() {
        var c models.Comment
        if err := rows.Scan(&c.ID, &c.PostID, &c.Author, &c.Text); err != nil {
            http.Error(w, "Error reading comment data", http.StatusInternalServerError)
            return
        }
        comments = append(comments, c)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(comments)
}

func createCommentHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    var c models.Comment
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

func updateCommentHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    var c models.Comment
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
