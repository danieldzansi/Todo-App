package services

import "github.com/danieldzansi/todo-api/internal/repository"


type TodoService interface{}

type todoService struct {
	repo repository.TodoRepository
}

func NewTodoService(r repository.TodoRepository) TodoService {
	return &todoService{repo: r}
}
