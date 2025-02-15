# Go Small Web App

This project is a simple web application written in Go. It connects to a PostgreSQL database and exposes RESTful endpoints for managing users, posts, and comments.

## Features
- Connects to a PostgreSQL database (with example connection string: `postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable`)
- Creates `users`, `posts`, and `comments` tables if they don’t already exist
- Exposes REST endpoints to **GET**, **POST**, and **PUT** records in each table
- Returns responses in JSON format

## Prerequisites
1. **Go** (version 1.18 or higher recommended)
2. **PostgreSQL** installed and running
3. A valid PostgreSQL user with the correct credentials. The default connection in the code is:
   ```
   postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable
   ```
   If needed, adjust it in the `initDB()` function in `main.go`.

## Installation and Setup
1. Clone or download this repository.
2. Ensure the `go.mod` and `go.sum` files are present. 
3. Run:
   ```
   go mod tidy
   ```
   to install required dependencies (like `github.com/lib/pq`).
4. Start the application:
   ```
   go run main.go
   ```
5. Once the application is running, it will listen on `http://localhost:8080`.  
   The first time it runs, it will create the `users`, `posts`, and `comments` tables in your PostgreSQL database if they do not already exist.

## Usage
Below are the endpoints and example **Windows** `curl` commands for **GET**, **POST**, and **PUT** requests.  
> **Note**: On Windows, you typically have `curl.exe` available (either built-in or installed). Use the exact syntax with quotes as shown below.

---

### 1. Users

#### GET all users
```
curl.exe -X GET "http://localhost:8080/users"
```
Example response:
```json
[
  {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com"
  },
  ...
]
```

#### POST (create new user)
```
curl.exe -X POST "http://localhost:8080/users" ^
    -H "Content-Type: application/json" ^
    -d "{\"name\":\"John Doe\",\"email\":\"john@example.com\"}"
```
Example response:
```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com"
}
```

#### PUT (update user)
```
curl.exe -X PUT "http://localhost:8080/users" ^
    -H "Content-Type: application/json" ^
    -d "{\"id\":1,\"name\":\"John Updated\",\"email\":\"john.updated@example.com\"}"
```
Example response:
```json
{
  "id": 1,
  "name": "John Updated",
  "email": "john.updated@example.com"
}
```

---

### 2. Posts

#### GET all posts
```
curl.exe -X GET "http://localhost:8080/posts"
```
Example response:
```json
[
  {
    "id": 1,
    "user_id": 1,
    "title": "Sample Post",
    "content": "Hello, this is a post"
  },
  ...
]
```

#### POST (create new post)
```
curl.exe -X POST "http://localhost:8080/posts" ^
    -H "Content-Type: application/json" ^
    -d "{\"user_id\":1,\"title\":\"My First Post\",\"content\":\"This is a post content.\"}"
```
Example response:
```json
{
  "id": 1,
  "user_id": 1,
  "title": "My First Post",
  "content": "This is a post content."
}
```

#### PUT (update post)
```
curl.exe -X PUT "http://localhost:8080/posts" ^
    -H "Content-Type: application/json" ^
    -d "{\"id\":1,\"title\":\"Updated Title\",\"content\":\"Updated post content.\"}"
```
Example response:
```json
{
  "id": 1,
  "user_id": 1,
  "title": "Updated Title",
  "content": "Updated post content."
}
```

---

### 3. Comments

#### GET all comments
```
curl.exe -X GET "http://localhost:8080/comments"
```
Example response:
```json
[
  {
    "id": 1,
    "post_id": 1,
    "author": "Alice",
    "text": "Nice post!"
  },
  ...
]
```

#### POST (create new comment)
```
curl.exe -X POST "http://localhost:8080/comments" ^
    -H "Content-Type: application/json" ^
    -d "{\"post_id\":1,\"author\":\"Alice\",\"text\":\"Great post content!\"}"
```
Example response:
```json
{
  "id": 1,
  "post_id": 1,
  "author": "Alice",
  "text": "Great post content!"
}
```

#### PUT (update comment)
```
curl.exe -X PUT "http://localhost:8080/comments" ^
    -H "Content-Type: application/json" ^
    -d "{\"id\":1,\"author\":\"Alice Updated\",\"text\":\"This comment has been updated.\"}"
```
Example response:
```json
{
  "id": 1,
  "post_id": 1,
  "author": "Alice Updated",
  "text": "This comment has been updated."
}
```

---

## License
This project is open-source and available under the [MIT License](https://opensource.org/licenses/MIT).
