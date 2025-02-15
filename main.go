package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// Переменная для хранения соединения с БД
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
	connStr :=
}