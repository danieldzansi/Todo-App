package repository

import (
	"database/sql"
	"fmt"
	"time"

	models "github.com/danieldzansi/todo-api/internal/model"
	"github.com/google/uuid"
)

type TodoRepository interface {
	CreateTodo(todo *models.Todo) error
}

type todoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) CreateTodo(todo *models.Todo) error {
	query := `
	  INSERT INTO todos(id,title,description,completed,due_date,created_at,updated_at)
	  VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	now := time.Now()
	todo.ID = uuid.New()
	todo.CreatedAt = now
	todo.UpdatedAt = now

	_, err := r.db.Exec(query,
		todo.ID,
		todo.Title,
		todo.Description,
		todo.Completed,
		todo.DueDate,
		todo.CreatedAt,
		todo.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create todo: %w", err)
	}
	return nil
}
