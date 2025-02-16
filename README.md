# Small Web App

This repository provides a small web application example using Go (Golang). It demonstrates how to structure a Go project with separate packages for database connections, handlers, middleware, and models. It also includes a basic JWT-based authentication flow and simple CRUD operations for users, posts, and comments.

## Features

- **JWT authentication** for user login and protected routes  
- **CRUD** operations for:
  - Users  
  - Posts  
  - Comments  
- **Modular architecture** using Go packages and a recommended folder structure  
- **PostgreSQL** database integration and automatic table creation if they do not exist  

## Folder Structure

```bash
go-small-webapp/
├── cmd/
│   └── webapp/
│       └── main.go
├── internal/
│   ├── database/
│   │   └── db.go
│   ├── handlers/
│   │   ├── auth.go
│   │   ├── user_handlers.go
│   │   ├── post_handlers.go
│   │   └── comment_handlers.go
│   ├── middleware/
│   │   └── auth_middleware.go
│   └── models/
│       ├── user.go
│       ├── post.go
│       └── comment.go
├── go.mod
└── go.sum
```

**Brief explanation of each main folder/file:**

- **cmd/webapp/main.go**  
  - Entry point of the application. Initializes the database, registers routes, and starts the HTTP server.

- **internal/database/db.go**  
  - Responsible for database connection setup and creating initial tables if they do not exist.

- **internal/handlers/\***  
  - Contains HTTP handler functions grouped by feature (authentication, users, posts, comments). 
  - Each file handles routing logic (GET, POST, PUT, etc.) for a specific resource.

- **internal/middleware/auth_middleware.go**  
  - Provides middleware for parsing and verifying JWT tokens before protected routes.

- **internal/models/\***  
  - Defines Go structs (models) representing database entities (User, Post, Comment).

## Getting Started

### Prerequisites

- **Go** (v1.19 or higher recommended)  
- **PostgreSQL** server running locally (or accessible remotely)  

### Installation

1. **Clone** this repository:
   ```bash
   git clone https://github.com/your-username/go-small-webapp.git
   ```
2. **Navigate** to the project folder:
   ```bash
   cd go-small-webapp
   ```
3. **Configure** your database settings:
   - By default, the connection string is located in `internal/database/db.go`. Update it if necessary to match your PostgreSQL credentials:
     ```go
     connStr := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
     ```
   - Alternatively, you can load environment variables in a custom way and build your connection string, but that would require updating the code to read from environment variables.

4. **Install dependencies**:
   ```bash
   go mod tidy
   ```

### Running the Application

From the project root directory, run:
```bash
go run ./cmd/webapp
```
This will:

1. Connect to your PostgreSQL database.  
2. Automatically create tables (`users`, `posts`, `comments`) if they do not exist.  
3. Start a local HTTP server at `http://localhost:8080`.

You should see output similar to:
```
Connection established
Server started on http://localhost:8080
```

## Usage

### Authentication

- **Endpoint**: `POST /login`  
- **Body** (JSON):
  ```json
  {
    "email": "your-email@example.com",
    "password": "your-secret-password"
  }
  ```
- **Response**:
  ```json
  {
    "token": "<JWT_TOKEN>"
  }
  ```

Use this token in the `Authorization` header as `Bearer <JWT_TOKEN>` for protected routes (e.g., `/posts`, `/comments`).

### Users

- **Endpoint**: `GET /users`  
  - Returns a list of all users (no token required in this example).

- **Endpoint**: `POST /users`  
  - Creates a new user.  
  - **Body** (JSON):
    ```json
    {
      "name": "John Doe",
      "email": "john@example.com",
      "password": "secret"
    }
    ```

- **Endpoint**: `PUT /users`  
  - Updates an existing user.  
  - **Body** (JSON):
    ```json
    {
      "id": 1,
      "name": "John Updated",
      "email": "johnupdated@example.com"
    }
    ```

### Posts

Protected by **JWT**:

- **Endpoint**: `GET /posts`  
  - Returns all posts.  
  - Requires valid `Authorization: Bearer <token>` header.

- **Endpoint**: `POST /posts`  
  - Creates a new post.  
  - **Body** (JSON):
    ```json
    {
      "user_id": 1,
      "title": "My First Post",
      "content": "This is the body of the post."
    }
    ```

- **Endpoint**: `PUT /posts`  
  - Updates an existing post.  
  - **Body** (JSON):
    ```json
    {
      "id": 1,
      "title": "Updated Title",
      "content": "Updated content."
    }
    ```

### Comments

Also protected by **JWT**:

- **Endpoint**: `GET /comments`  
- **Endpoint**: `POST /comments`  
- **Endpoint**: `PUT /comments`  

Example **POST** request body:
```json
{
  "post_id": 1,
  "author": "Some Author",
  "text": "This is a comment."
}
```

## Customization

- Modify the connection string in `internal/database/db.go` or replace it with environment variable logic to suit your environment.  
- Adjust the table creation in `CreateTables` if you need additional fields or different schema.  
- Update the routes or split them further if your project grows.

## License

This project is distributed under the MIT License. Feel free to use it as a template or a starting point for your own Go web applications.