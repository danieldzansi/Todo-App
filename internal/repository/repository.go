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
	GetAllTodos() ([]models.Todo, error)
	GetTodoByID(id uuid.UUID) (*models.Todo, error)
	UpdateTodo(id uuid.UUID, req *models.UpdateTodoRequest) (*models.Todo, error)
	DeleteTodo(id uuid.UUID) error
	ToggleTodoComplete(id uuid.UUID) (*models.Todo, error)
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

func (r *todoRepository) GetAllTodos() ([]models.Todo, error) {
	query := `
	  SELECT id, title, description, completed, due_date, created_at, updated_at
	  FROM todos
	  ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query todos: %w", err)
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var t models.Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Completed, &t.DueDate, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan todo: %w", err)
		}
		todos = append(todos, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}
	return todos, nil
}

func (r *todoRepository) GetTodoByID(id uuid.UUID) (*models.Todo, error) {
	query := `
	  SELECT id, title, description, completed, due_date, created_at, updated_at
	  FROM todos
	  WHERE id = $1
	`
	var t models.Todo
	err := r.db.QueryRow(query, id).Scan(&t.ID, &t.Title, &t.Description, &t.Completed, &t.DueDate, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrTodoNotFound
		}
		return nil, fmt.Errorf("failed to get todo by id: %w", err)
	}
	return &t, nil
}

func (r *todoRepository) UpdateTodo(id uuid.UUID, req *models.UpdateTodoRequest) (*models.Todo, error) {
	existing, err := r.GetTodoByID(id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		existing.Title = *req.Title
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.DueDate != nil {
		existing.DueDate = req.DueDate
	}
	now := time.Now()
	existing.UpdatedAt = now

	query := `
	  UPDATE todos
	  SET title = $1, description = $2, due_date = $3, updated_at = $4
	  WHERE id = $5
	`
	res, err := r.db.Exec(query, existing.Title, existing.Description, existing.DueDate, existing.UpdatedAt, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update todo: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}
	if affected == 0 {
		return nil, models.ErrTodoNotFound
	}

	return existing, nil
}

func (r *todoRepository) DeleteTodo(id uuid.UUID) error {
	query := `DELETE FROM todos WHERE id = $1`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if affected == 0 {
		return models.ErrTodoNotFound
	}
	return nil
}

func (r *todoRepository) ToggleTodoComplete(id uuid.UUID) (*models.Todo, error) {
	existing, err := r.GetTodoByID(id)
	if err != nil {
		return nil, err
	}
	existing.Completed = !existing.Completed
	existing.UpdatedAt = time.Now()

	query := `
	  UPDATE todos
	  SET completed = $1, updated_at = $2
	  WHERE id = $3
	`
	res, err := r.db.Exec(query, existing.Completed, existing.UpdatedAt, id)
	if err != nil {
		return nil, fmt.Errorf("failed to toggle todo: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}
	if affected == 0 {
		return nil, models.ErrTodoNotFound
	}
	return existing, nil
}
