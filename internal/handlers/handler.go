package handlers

import (
	"bytes"
	"encoding/json"
	"io"
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

// DummyJSON base URL
const dummyURL = "https://dummyjson.com"

// Reusable helper to call DummyJSON API
func fetchDummyJSON(method, endpoint string, body interface{}) ([]byte, error) {
	url := dummyURL + endpoint

	var req *http.Request
	var err error

	if body != nil {
		jsonBody, _ := json.Marshal(body)
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

//
// ====== HANDLERS (switch between Local DB & DummyJSON) ======
//

func (h *TodoHandler) GetAllTodos(c *gin.Context) {
	if c.Query("source") == "online" {
		data, err := fetchDummyJSON("GET", "/todos", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Data(http.StatusOK, "application/json", data)
		return
	}

	// Local DB
	todos, err := h.svc.GetAllTodos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"todos": todos})
}

func (h *TodoHandler) GetTodoByID(c *gin.Context) {
	if c.Query("source") == "online" {
		id := c.Param("id")
		data, err := fetchDummyJSON("GET", "/todos/"+id, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Data(http.StatusOK, "application/json", data)
		return
	}

	// Local DB
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
	if c.Query("source") == "online" {
		var newTodo map[string]interface{}
		if err := c.ShouldBindJSON(&newTodo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		data, err := fetchDummyJSON("POST", "/todos/add", newTodo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Data(http.StatusCreated, "application/json", data)
		return
	}

	// Local DB
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
	if c.Query("source") == "online" {
		id := c.Param("id")
		var updateTodo map[string]interface{}
		if err := c.ShouldBindJSON(&updateTodo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		data, err := fetchDummyJSON("PUT", "/todos/"+id, updateTodo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Data(http.StatusOK, "application/json", data)
		return
	}

	// Local DB
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
	if c.Query("source") == "online" {
		id := c.Param("id")
		data, err := fetchDummyJSON("DELETE", "/todos/"+id, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Data(http.StatusOK, "application/json", data)
		return
	}

	// Local DB
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
	// NOTE: DummyJSON doesn’t support PATCH complete toggle directly
	if c.Query("source") == "online" {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "toggle complete not supported on DummyJSON"})
		return
	}

	// Local DB
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
