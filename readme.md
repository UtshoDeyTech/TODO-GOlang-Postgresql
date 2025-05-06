# TODO List Backend in Go with PostgreSQL

This is a simple RESTful API for a TODO list application built with Go and PostgreSQL. The project includes CRUD operations for managing todo items and is containerized using Docker.

## Project Structure

```
todo-api/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── handlers/
│   │   └── todo_handler.go
│   ├── models/
│   │   └── todo.go
│   ├── repository/
│   │   └── todo_repository.go
│   └── router/
│       └── router.go
├── migrations/
│   └── init.sql
├── .env.example
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── README.md
```

## Code Files

### `cmd/api/main.go`

```go
package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourusername/todo-api/internal/config"
	"github.com/yourusername/todo-api/internal/router"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize router
	r := router.SetupRouter(cfg)

	// Configure the HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
}
```

### `internal/config/config.go`

```go
package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	DB       *sql.DB
	DBConfig DBConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	// Set default values
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Database configuration
	dbConfig := DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "todo_db"),
	}

	// Connect to database
	db, err := connectDB(dbConfig)
	if err != nil {
		return nil, err
	}

	return &Config{
		Port:     port,
		DB:       db,
		DBConfig: dbConfig,
	}, nil
}

// connectDB establishes a connection to the database
func connectDB(config DBConfig) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// getEnv gets the value of an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
```

### `internal/models/todo.go`

```go
package models

import "time"

// Todo represents a todo item
type Todo struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// CreateTodoRequest represents the request payload for creating a todo
type CreateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// UpdateTodoRequest represents the request payload for updating a todo
type UpdateTodoRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Completed   *bool   `json:"completed,omitempty"`
}
```

### `internal/repository/todo_repository.go`

```go
package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/yourusername/todo-api/internal/models"
)

// TodoRepository handles database operations for todos
type TodoRepository struct {
	db *sql.DB
}

// NewTodoRepository creates a new TodoRepository
func NewTodoRepository(db *sql.DB) *TodoRepository {
	return &TodoRepository{
		db: db,
	}
}

// Create adds a new todo to the database
func (r *TodoRepository) Create(todo *models.CreateTodoRequest) (*models.Todo, error) {
	query := `
		INSERT INTO todos (title, description, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id, title, description, completed, created_at, updated_at, completed_at
	`

	var newTodo models.Todo
	var completedAt sql.NullTime

	err := r.db.QueryRow(
		query,
		todo.Title,
		todo.Description,
	).Scan(
		&newTodo.ID,
		&newTodo.Title,
		&newTodo.Description,
		&newTodo.Completed,
		&newTodo.CreatedAt,
		&newTodo.UpdatedAt,
		&completedAt,
	)

	if err != nil {
		return nil, err
	}

	if completedAt.Valid {
		newTodo.CompletedAt = &completedAt.Time
	}

	return &newTodo, nil
}

// GetAll retrieves all todos from the database
func (r *TodoRepository) GetAll() ([]*models.Todo, error) {
	query := `
		SELECT id, title, description, completed, created_at, updated_at, completed_at
		FROM todos
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*models.Todo

	for rows.Next() {
		var todo models.Todo
		var completedAt sql.NullTime

		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
			&completedAt,
		)

		if err != nil {
			return nil, err
		}

		if completedAt.Valid {
			todo.CompletedAt = &completedAt.Time
		}

		todos = append(todos, &todo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

// GetByID retrieves a todo by ID
func (r *TodoRepository) GetByID(id int64) (*models.Todo, error) {
	query := `
		SELECT id, title, description, completed, created_at, updated_at, completed_at
		FROM todos
		WHERE id = $1
	`

	var todo models.Todo
	var completedAt sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&completedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Todo not found
		}
		return nil, err
	}

	if completedAt.Valid {
		todo.CompletedAt = &completedAt.Time
	}

	return &todo, nil
}

// Update updates a todo in the database
func (r *TodoRepository) Update(id int64, todo *models.UpdateTodoRequest) (*models.Todo, error) {
	// First, get the current todo
	currentTodo, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	if currentTodo == nil {
		return nil, nil // Todo not found
	}

	// Prepare update values
	title := currentTodo.Title
	if todo.Title != nil {
		title = *todo.Title
	}

	description := currentTodo.Description
	if todo.Description != nil {
		description = *todo.Description
	}

	completed := currentTodo.Completed
	var completedAt *time.Time = currentTodo.CompletedAt

	if todo.Completed != nil && *todo.Completed != completed {
		completed = *todo.Completed
		if completed {
			now := time.Now()
			completedAt = &now
		} else {
			completedAt = nil
		}
	}

	// Update in database
	query := `
		UPDATE todos
		SET title = $1, description = $2, completed = $3, completed_at = $4, updated_at = NOW()
		WHERE id = $5
		RETURNING id, title, description, completed, created_at, updated_at, completed_at
	`

	var updatedTodo models.Todo
	var nullCompletedAt sql.NullTime
	if completedAt != nil {
		nullCompletedAt = sql.NullTime{Time: *completedAt, Valid: true}
	}

	err = r.db.QueryRow(
		query,
		title,
		description,
		completed,
		nullCompletedAt,
		id,
	).Scan(
		&updatedTodo.ID,
		&updatedTodo.Title,
		&updatedTodo.Description,
		&updatedTodo.Completed,
		&updatedTodo.CreatedAt,
		&updatedTodo.UpdatedAt,
		&nullCompletedAt,
	)

	if err != nil {
		return nil, err
	}

	if nullCompletedAt.Valid {
		updatedTodo.CompletedAt = &nullCompletedAt.Time
	}

	return &updatedTodo, nil
}

// Delete removes a todo from the database
func (r *TodoRepository) Delete(id int64) error {
	query := `DELETE FROM todos WHERE id = $1`

	_, err := r.db.Exec(query, id)
	return err
}
```

### `internal/handlers/todo_handler.go`

```go
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yourusername/todo-api/internal/models"
	"github.com/yourusername/todo-api/internal/repository"
)

// TodoHandler handles HTTP requests for todo operations
type TodoHandler struct {
	repo *repository.TodoRepository
}

// NewTodoHandler creates a new TodoHandler
func NewTodoHandler(repo *repository.TodoRepository) *TodoHandler {
	return &TodoHandler{
		repo: repo,
	}
}

// GetAllTodos handles GET /todos
func (h *TodoHandler) GetAllTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, todos)
}

// GetTodo handles GET /todos/{id}
func (h *TodoHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	todo, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if todo == nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	respondWithJSON(w, http.StatusOK, todo)
}

// CreateTodo handles POST /todos
func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate request
	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	todo, err := h.repo.Create(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusCreated, todo)
}

// UpdateTodo handles PUT /todos/{id}
func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	todo, err := h.repo.Update(id, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if todo == nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	respondWithJSON(w, http.StatusOK, todo)
}

// DeleteTodo handles DELETE /todos/{id}
func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// respondWithJSON writes the response as JSON
func respondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
```

### `internal/router/router.go`

```go
package router

import (
	"github.com/gorilla/mux"
	"github.com/yourusername/todo-api/internal/config"
	"github.com/yourusername/todo-api/internal/handlers"
	"github.com/yourusername/todo-api/internal/repository"
)

// SetupRouter configures the HTTP router
func SetupRouter(cfg *config.Config) *mux.Router {
	r := mux.NewRouter()

	// Initialize repositories
	todoRepo := repository.NewTodoRepository(cfg.DB)

	// Initialize handlers
	todoHandler := handlers.NewTodoHandler(todoRepo)

	// Define API routes
	api := r.PathPrefix("/api/v1").Subrouter()

	// Todo routes
	api.HandleFunc("/todos", todoHandler.GetAllTodos).Methods("GET")
	api.HandleFunc("/todos/{id:[0-9]+}", todoHandler.GetTodo).Methods("GET")
	api.HandleFunc("/todos", todoHandler.CreateTodo).Methods("POST")
	api.HandleFunc("/todos/{id:[0-9]+}", todoHandler.UpdateTodo).Methods("PUT")
	api.HandleFunc("/todos/{id:[0-9]+}", todoHandler.DeleteTodo).Methods("DELETE")

	return r
}
```

### `migrations/init.sql`

```sql
CREATE TABLE IF NOT EXISTS todos (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP
);

-- Add some sample data
INSERT INTO todos (title, description, completed, created_at, updated_at, completed_at)
VALUES
    ('Learn Go', 'Study Go programming language basics', true, NOW(), NOW(), NOW()),
    ('Build a REST API', 'Create a Todo API using Go and PostgreSQL', false, NOW(), NOW(), NULL),
    ('Learn Docker', 'Learn how to containerize applications', false, NOW(), NOW(), NULL);
```

### `Dockerfile`

```dockerfile
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./
# Copy source code
COPY . .

# Build the application
RUN go build -o main ./cmd/api/main.go

# Create a minimal image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/.env* ./

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
```

### `docker-compose.yml`

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      - db
    environment:
      - DB_HOST=db
      - DB_USER=postgres
      - DB_PASSWORD=postgres
      - DB_NAME=todo_db
      - DB_PORT=5432
      - PORT=8080
    networks:
      - todo-network
    restart: unless-stopped

  db:
    image: postgres:14-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: todo_db
    ports:
      - "5432:5432"
    volumes:
      - postgres-data:/var/lib/postgresql/data
      - ./migrations/init.sql:/docker-entrypoint-initdb.d/init.sql
    networks:
      - todo-network
    restart: unless-stopped

networks:
  todo-network:
    driver: bridge

volumes:
  postgres-data:
```

### `go.mod`

```go
module github.com/yourusername/todo-api

go 1.20

require (
	github.com/gorilla/mux v1.8.0
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
)
```

### `.env.example`

```
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=todo_db
```

### `.gitignore`

```
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Environment variables
.env

# Test binary, built with `go test -c`
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below if you want to check in your vendor directory)
vendor/

# Go workspace file
go.work

# Local development artifacts
tmp/
dist/
```

## API Endpoints

- `GET /api/v1/todos` - Get all todos
- `GET /api/v1/todos/{id}` - Get a specific todo
- `POST /api/v1/todos` - Create a new todo
- `PUT /api/v1/todos/{id}` - Update a todo
- `DELETE /api/v1/todos/{id}` - Delete a todo

## Building and Running

1. Clone the repository
2. Make sure Docker is installed
3. Run `docker-compose up --build`
4. The API will be available at `http://localhost:8080/api/v1/todos`

## Example API Requests

### Create a Todo

```
POST /api/v1/todos
Content-Type: application/json

{
  "title": "Buy groceries",
  "description": "Milk, eggs, bread, and cheese"
}
```

### Update a Todo

```
PUT /api/v1/todos/1
Content-Type: application/json

{
  "title": "Buy organic groceries",
  "completed": true
}
```

### Get All Todos

```
GET /api/v1/todos
```

### Get a Todo

```
GET /api/v1/todos/1
```

### Delete a Todo

```
DELETE /api/v1/todos/1
```