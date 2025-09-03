package handlers

import (
	"net/http"

	models "github.com/danieldzansi/todo-api/internal/model"
	"github.com/danieldzansi/todo-api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TodoHandler struct {
	svc services.TodoService
}

func NewTodoHandler(s services.TodoService) *TodoHandler {
	return &TodoHandler{svc: s}
}

func (h *TodoHandler) GetAllTodos(c *gin.Context) {
	todos, err := h.svc.GetAllTodos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"todos": todos})
}

func (h *TodoHandler) GetTodoByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	todo, err := h.svc.GetTodoByID(id)
	if err != nil {
		if err == models.ErrTodoNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"todo": todo})
}

func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var req models.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	todo, err := h.svc.CreateTodo(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"todo": todo})
}

func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req models.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	todo, err := h.svc.UpdateTodo(id, &req)
	if err != nil {
		if err == models.ErrTodoNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"todo": todo})
}

func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.svc.DeleteTodo(id); err != nil {
		if err == models.ErrTodoNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *TodoHandler) ToggleTodoComplete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	todo, err := h.svc.ToggleTodoComplete(id)
	if err != nil {
		if err == models.ErrTodoNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"todo": todo})
}
