package main

import (
    "fmt"
    "log"
    "net/http"

    "small_web_app/internal/database"
    "small_web_app/internal/handlers"
    "small_web_app/internal/middleware"
)

func main() {
    // Инициализируем базу
    db, err := database.InitDB()
    if err != nil {
        log.Fatalf("Ошибка инициализации БД: %v", err)
    }
    defer db.Close()

    // Создаём таблицы (если нужно)
    if err := database.CreateTables(db); err != nil {
        log.Printf("Ошибка при создании таблиц: %v", err)
    }

    // Регистрируем все маршруты
    mux := http.NewServeMux()

    // Маршрут для логина (не требует middleware)
    mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        handlers.LoginHandler(db, w, r)
    })

    // Группа маршрутов /users
    mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
        handlers.UserHandlers(db, w, r)
    })

    // Пример – для защищённых маршрутов используем middleware
    mux.Handle("/posts", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        handlers.PostHandlers(db, w, r)
    })))

    mux.Handle("/comments", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        handlers.CommentHandlers(db, w, r)
    })))

    fmt.Println("Сервер запущен на http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}
