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
    // DB initialization
    db, err := database.InitDB()
    if err != nil {
        log.Fatalf("Error initializing the database: %v", err)
    }
    defer db.Close()

    // Creating tables
    if err := database.CreateTables(db); err != nil {
        log.Printf("Error creating tables: %v", err)
    }

    // Register all routes
    mux := http.NewServeMux()

    // Login route (does not require middleware)
    mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        handlers.LoginHandler(db, w, r)
    })

    // User routes group
    mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
        handlers.UserHandlers(db, w, r)
    })

    // Example – use middleware for protected routes
    mux.Handle("/posts", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        handlers.PostHandlers(db, w, r)
    })))

    mux.Handle("/comments", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        handlers.CommentHandlers(db, w, r)
    })))

    fmt.Println("Server started at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}
