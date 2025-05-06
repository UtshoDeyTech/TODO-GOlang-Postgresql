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
