package handlers

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"

    "golang.org/x/crypto/bcrypt"

    "small_web_app/internal/models"
)

// Один общий router-хендлер на /users
// Выбираем метод запроса и вызываем нужную функцию
func UserHandlers(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        getUsersHandler(db, w, r)
    case http.MethodPost:
        createUserHandler(db, w, r)
    case http.MethodPut:
        updateUserHandler(db, w, r)
    default:
        http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
    }
}

func getUsersHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT id, name, email FROM users")
    if err != nil {
        http.Error(w, "Ошибка при запросе пользователей", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var users []models.User
    for rows.Next() {
        var u models.User
        if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
            http.Error(w, "Ошибка чтения данных пользователя", http.StatusInternalServerError)
            return
        }
        users = append(users, u)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

func createUserHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    var newUser models.User
    if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
        http.Error(w, "Некорректные данные пользователя", http.StatusBadRequest)
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Ошибка при хешировании пароля", http.StatusInternalServerError)
        return
    }

    query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id`
    err = db.QueryRow(query, newUser.Name, newUser.Email, string(hashedPassword)).Scan(&newUser.ID)
    if err != nil {
        http.Error(w, "Ошибка при создании пользователя", http.StatusInternalServerError)
        return
    }

    // чтобы не возвращать hash в ответе
    newUser.Password = ""

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(newUser)
}

func updateUserHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    var u models.User
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
