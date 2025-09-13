package main

import (
	"log"
	"os"

	database "github.com/danieldzansi/todo-api/internal/database"
	"github.com/danieldzansi/todo-api/internal/handlers"
	"github.com/danieldzansi/todo-api/internal/repository"
	"github.com/danieldzansi/todo-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	conn, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer conn.Close()

	todoRepo := repository.NewTodoRepository(conn)
	todoService := services.NewTodoService(todoRepo)

	userRepo := repository.NewUserRepository(conn)
	authService := services.NewAuthService(userRepo)

	todoHandler := handlers.NewTodoHandler(todoService)
	userHandler := handlers.NewUserHandler(authService)

	gin.SetMode(os.Getenv("GIN_MODE"))

	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	api := router.Group("/api/v1")
	{

		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"message": "Todo API is running",
			})
		})

		todos := api.Group("/todos")
		{
			// protect todos with JWT
			todos.Use(handlers.AuthMiddleware())
			todos.GET("/", todoHandler.GetAllTodos)
			todos.GET("/:id", todoHandler.GetTodoByID)
			todos.POST("/", todoHandler.CreateTodo)
			todos.PUT("/:id", todoHandler.UpdateTodo)
			todos.DELETE("/:id", todoHandler.DeleteTodo)
			todos.PATCH("/:id/complete", todoHandler.ToggleTodoComplete)
		}

		users := api.Group("/users")
		{
			users.POST("/signup", userHandler.Signup)
			users.POST("/login", userHandler.Login)
			users.GET("/:id", userHandler.GetUserByID)
			users.GET("/", userHandler.GetAllUsers)
		}
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
