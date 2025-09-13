package services

import (
	"errors"
	"os"
	"time"

	models "github.com/danieldzansi/todo-api/internal/model"
	"github.com/danieldzansi/todo-api/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TodoService interface {
	CreateTodo(req *models.CreateTodoRequest) (*models.Todo, error)
	GetAllTodos() ([]models.Todo, error)
	GetTodoByID(id uuid.UUID) (*models.Todo, error)
	UpdateTodo(id uuid.UUID, req *models.UpdateTodoRequest) (*models.Todo, error)
	DeleteTodo(id uuid.UUID) error
	ToggleTodoComplete(id uuid.UUID) (*models.Todo, error)
	CreateTodoForUser(userID uuid.UUID, req *models.CreateTodoRequest) (*models.Todo, error)
	GetAllTodosByUser(userID uuid.UUID) ([]models.Todo, error)
	GetTodoByIDForUser(userID uuid.UUID, id uuid.UUID) (*models.Todo, error)
	UpdateTodoForUser(userID uuid.UUID, id uuid.UUID, req *models.UpdateTodoRequest) (*models.Todo, error)
	DeleteTodoForUser(userID uuid.UUID, id uuid.UUID) error
	ToggleTodoCompleteForUser(userID uuid.UUID, id uuid.UUID) (*models.Todo, error)
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
func (s *todoService) CreateTodoForUser(userID uuid.UUID, req *models.CreateTodoRequest) (*models.Todo, error) {
	todo := &models.Todo{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		Completed:   false,
		UserID:      userID,
	}
	if err := s.repo.CreateTodo(todo); err != nil {
		return nil, err
	}
	return todo, nil
}

type AuthService interface {
	Signup(req *models.SignupRequest) (*models.User, error)
	Login(req *models.LoginRequest) (*models.LoginResponse, error)
	GetUserByID(id uuid.UUID) (*models.User, error)
	GetAllUsers() ([]models.User, error)
}

type AuthServiceImpl struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &AuthServiceImpl{repo: repo}
}

func (s *AuthServiceImpl) Signup(req *models.SignupRequest) (*models.User, error) {
	if existing, err := s.repo.GetUserByEmail(req.Email); err == nil && existing != nil {
		return nil, models.ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := &models.User{
		ID:        uuid.New(),
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashedPassword),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.CreateUser(*user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthServiceImpl) GetUserByID(id uuid.UUID) (*models.User, error) {
	return s.repo.GetUserByID(id)
}
func (s *AuthServiceImpl) GetAllUsers() ([]models.User, error) {
	return s.repo.GetAllUsers()
}

func (s *AuthServiceImpl) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return nil, errors.New("invalid credentials")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-change-me"
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID.String(),
		"email": user.Email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	})
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}
	safeUser := *user
	safeUser.Password = ""

	return &models.LoginResponse{Token: signed, User: safeUser}, nil
}

func (s *todoService) GetAllTodos() ([]models.Todo, error) {
	return s.repo.GetAllTodos()
}
func (s *todoService) GetAllTodosByUser(userID uuid.UUID) ([]models.Todo, error) {
	return s.repo.GetAllTodosByUser(userID)
}

func (s *todoService) GetTodoByID(id uuid.UUID) (*models.Todo, error) {
	return s.repo.GetTodoByID(id)
}
func (s *todoService) GetTodoByIDForUser(userID uuid.UUID, id uuid.UUID) (*models.Todo, error) {
	return s.repo.GetTodoByIDForUser(userID, id)
}

func (s *todoService) UpdateTodo(id uuid.UUID, req *models.UpdateTodoRequest) (*models.Todo, error) {
	return s.repo.UpdateTodo(id, req)
}
func (s *todoService) UpdateTodoForUser(userID uuid.UUID, id uuid.UUID, req *models.UpdateTodoRequest) (*models.Todo, error) {
	return s.repo.UpdateTodoForUser(userID, id, req)
}

func (s *todoService) DeleteTodo(id uuid.UUID) error {
	return s.repo.DeleteTodo(id)
}
func (s *todoService) DeleteTodoForUser(userID uuid.UUID, id uuid.UUID) error {
	return s.repo.DeleteTodoForUser(userID, id)
}

func (s *todoService) ToggleTodoComplete(id uuid.UUID) (*models.Todo, error) {
	return s.repo.ToggleTodoComplete(id)
}
func (s *todoService) ToggleTodoCompleteForUser(userID uuid.UUID, id uuid.UUID) (*models.Todo, error) {
	return s.repo.ToggleTodoCompleteForUser(userID, id)
}
