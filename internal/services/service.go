package services

import (
	models "github.com/danieldzansi/todo-api/internal/model"
	"github.com/danieldzansi/todo-api/internal/repository"
	"github.com/google/uuid"
)

type TodoService interface {
	CreateTodo(req *models.CreateTodoRequest) (*models.Todo, error)
	GetAllTodos() ([]models.Todo, error)
	GetTodoByID(id uuid.UUID) (*models.Todo, error)
	UpdateTodo(id uuid.UUID, req *models.UpdateTodoRequest) (*models.Todo, error)
	DeleteTodo(id uuid.UUID) error
	ToggleTodoComplete(id uuid.UUID) (*models.Todo, error)
}

type todoService struct {
	repo repository.TodoRepository
}

func NewTodoService(r repository.TodoRepository) TodoService {
	return &todoService{repo: r}
}
func (s *todoService) CreateTodo(req *models.CreateTodoRequest) (*models.Todo, error) {
	todo := &models.Todo{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		Completed:   false,
	}

	err := s.repo.CreateTodo(todo)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *todoService) GetAllTodos() ([]models.Todo, error) {
	return s.repo.GetAllTodos()
}

func (s *todoService) GetTodoByID(id uuid.UUID) (*models.Todo, error) {
	return s.repo.GetTodoByID(id)
}

func (s *todoService) UpdateTodo(id uuid.UUID, req *models.UpdateTodoRequest) (*models.Todo, error) {
	return s.repo.UpdateTodo(id, req)
}

func (s *todoService) DeleteTodo(id uuid.UUID) error {
	return s.repo.DeleteTodo(id)
}

func (s *todoService) ToggleTodoComplete(id uuid.UUID) (*models.Todo, error) {
	return s.repo.ToggleTodoComplete(id)
}
