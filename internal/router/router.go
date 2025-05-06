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
