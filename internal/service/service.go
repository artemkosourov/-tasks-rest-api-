package service

import (
	"context"
	"tasks-rest-api/internal/client"
	"tasks-rest-api/internal/dto"
	"tasks-rest-api/internal/kafka"
	"tasks-rest-api/internal/model"
	"tasks-rest-api/internal/repository"
	"tasks-rest-api/internal/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type TaskService struct {
	repo          *repository.PostgresRepo
	kafkaProducer *kafka.Producer
	usersClient   *client.UsersAPIClient
}

func NewTaskService(repo *repository.PostgresRepo, kafkaProducer *kafka.Producer, userClient *client.UsersAPIClient) *TaskService {
	return &TaskService{
		repo:          repo,
		kafkaProducer: kafkaProducer,
		usersClient:   userClient,
	}
}

func (s *TaskService) RegisterUser(req dto.RegisterRequest) error {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := &model.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "user",
	}
	if err := s.repo.CreateUser(user); err != nil {
		return err
	}

	// Produce в кафку
	if s.kafkaProducer != nil {
		go func() {
			event := kafka.UserEvent{
				EventType: "REGISTERED",
				UserID:    user.ID,
				Email:     user.Email,
				Timestamp: time.Now(),
			}
			s.kafkaProducer.PublishUserEvent(event)
		}()
	}

	return nil
}

func (s *TaskService) AuthenticateUser(email string, password string) (*dto.AuthResponse, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}
	accessToken, _ := utils.GenerateAccessToken(user.ID, user.Email, user.Role)
	refreshToken, _ := utils.GenerateRefreshToken(user.ID, user.Email, user.Role)
	// Produce в кафку
	if s.kafkaProducer != nil {
		go func() {
			event := kafka.UserEvent{
				EventType: "LOGIN",
				UserID:    user.ID,
				Email:     user.Email,
				Timestamp: time.Now(),
			}
			s.kafkaProducer.PublishUserEvent(event)
		}()
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *TaskService) CreateTask(userID int, req dto.TaskRequest) error {
	task := &model.Task{
		Title:  req.Title,
		UserID: userID,
	}
	return s.repo.CreateTask(task)
}

func (s *TaskService) GetTasks(userID int) ([]model.Task, error) {
	return s.repo.GetTasksByUserId(userID)
}

func (s *TaskService) GetExternalUsers(ctx context.Context) ([]client.User, error) {
	return s.usersClient.GetAllUsers(ctx)
}
