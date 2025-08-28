package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrTodoNotFound = errors.New("todo not found")

type Todo struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Title       string     `json:"title" db:"title" binding:"required"`
	Description string     `json:"description" db:"description"`
	Completed   bool       `json:"completed" db:"completed"`
	Priority    int        `json:"priority" db:"priority"`
	DueDate     *time.Time `json:"due_date,omitempty" db:"due_date"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

type CreateTodoRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Priority    int        `json:"priority"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

type UpdateTodoRequest struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	Priority    *int       `json:"priority,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

type TodoResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
