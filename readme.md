# Todo List API

A simple and efficient RESTful API for managing todo items, built with Go and PostgreSQL.

## Overview

This Todo List API provides a complete backend solution for managing tasks with full CRUD (Create, Read, Update, Delete) operations. It's built using Go's standard libraries and popular packages like Gorilla Mux for routing, and uses PostgreSQL for data persistence.

## Features

- **Create**: Add new todo items with title and description
- **Read**: Get all todos or retrieve a specific todo by ID
- **Update**: Modify existing todo items (title, description, completion status)
- **Delete**: Remove todo items
- **Timestamps**: Automatic tracking of creation, update, and completion times

## Tech Stack

- **Go**: Backend language
- **PostgreSQL**: Database for persistent storage
- **Docker**: Container for PostgreSQL (optional)
- **Gorilla Mux**: HTTP router and URL matcher

## Project Structure

```

todo-api/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go               # Configuration management
│   ├── handlers/
│   │   └── todo\_handler.go         # HTTP request handlers
│   ├── models/
│   │   └── todo.go                 # Data models
│   ├── repository/
│   │   └── todo\_repository.go      # Database operations
│   └── router/
│       └── router.go               # API routes configuration
├── migrations/
│   └── init.sql                    # Database initialization script
├── .env.example                    # Example environment variables
├── docker-compose.yml              # Docker configuration for PostgreSQL
├── go.mod                          # Go module definition
├── go.sum                          # Go module checksums
└── README.md                       # Project documentation

````

## API Endpoints

| Method | Endpoint              | Description           | Request Body                                | Response                |
|--------|----------------------|----------------------|---------------------------------------------|-------------------------|
| GET    | /api/v1/todos        | Get all todos        | -                                           | Array of todo objects   |
| GET    | /api/v1/todos/{id}   | Get todo by ID       | -                                           | Single todo object      |
| POST   | /api/v1/todos        | Create a new todo    | `{"title": "...", "description": "..."}`    | Created todo object     |
| PUT    | /api/v1/todos/{id}   | Update a todo        | `{"title": "...", "completed": true}`       | Updated todo object     |
| DELETE | /api/v1/todos/{id}   | Delete a todo        | -                                           | No content              |

## Getting Started

### Prerequisites

- Go 1.20 or higher
- PostgreSQL (or Docker for containerized PostgreSQL)

### Setup and Installation

1. **Clone the repository**

```bash
git clone https://github.com/yourusername/todo-api.git
cd todo-api
````

2. **Set up the database**

Option A: Using Docker (recommended):

```bash
docker compose up -d
```

Option B: Using existing PostgreSQL installation:

* Create a database named `todo_db`
* Run the SQL script in `migrations/init.sql`

3. **Configure environment variables**

```bash
cp .env.example .env
# Edit .env with your database credentials if needed
```

4. **Build and run the application**

```bash
go mod tidy
go build -o main ./cmd/api/main.go
./main
```

The API will be available at [http://localhost:8080/api/v1/todos](http://localhost:8080/api/v1/todos)

## Example Usage

### Create a Todo

```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "Buy groceries", "description": "Milk, eggs, bread, and cheese"}'
```

### Get All Todos

```bash
curl http://localhost:8080/api/v1/todos
```

### Update a Todo

```bash
curl -X PUT http://localhost:8080/api/v1/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Buy organic groceries", "completed": true}'
```

### Delete a Todo

```bash
curl -X DELETE http://localhost:8080/api/v1/todos/1
```

## Development

### Running Tests

```bash
go test ./...
```

### Adding New Features

1. Create appropriate models in the models package
2. Implement repository methods in the repository package
3. Add handlers in the handlers package
4. Register new routes in the router package

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

* [Gorilla Mux](https://github.com/gorilla/mux) for HTTP routing
* [godotenv](https://github.com/joho/godotenv) for environment variable management
* [pq](https://github.com/lib/pq) for PostgreSQL driver



