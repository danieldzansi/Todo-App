package handlers

import (
	"net/http"

	"github.com/danieldzansi/todo-api/internal/services"
	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	svc services.TodoService
}

func NewTodoHandler(s services.TodoService) *TodoHandler {
	return &TodoHandler{svc: s}
}

func (h *TodoHandler) GetAllTodos(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"todos": []any{}})
}

func (h *TodoHandler) GetTodoByID(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"todo": gin.H{"id": c.Param("id")}})
}

func (h *TodoHandler) CreateTodo(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"status": "created"})
}

func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	c.JSON(http.StatusNoContent, nil)
}

func (h *TodoHandler) ToggleTodoComplete(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "toggled"})
}
