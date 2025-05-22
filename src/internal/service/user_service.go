package service

import (
	"smart-contract-automation/src/internal/model"
	"smart-contract-automation/src/internal/repository"
	"smart-contract-automation/src/pkg/logger"

	"go.uber.org/zap"
)

type UserService interface {
	CreateUser(user *model.User, password string) error
	ListAdmins() ([]model.User, error)
	GetUser(id string) (*model.User, error)
	DeleteUser(id string) error
	HasAnyUsers() (bool, error)
}

type userService struct {
	userRepo repository.UserRepository
	auth     AuthService
}

func NewUserService(userRepo repository.UserRepository, auth AuthService) UserService {
	return &userService{
		userRepo: userRepo,
		auth:     auth,
	}
}

func (s *userService) CreateUser(user *model.User, password string) error {
	// Hash the password
	hashedPassword, err := s.auth.HashPassword(password)
	if err != nil {
		logger.Error("Failed to hash password", zap.Error(err))
		return err
	}

	user.Password = hashedPassword
	return s.userRepo.Create(user)
}

func (s *userService) ListAdmins() ([]model.User, error) {
	return s.userRepo.FindByRoles([]model.Role{model.RoleAdmin, model.RoleSuperAdmin})
}

func (s *userService) GetUser(id string) (*model.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *userService) DeleteUser(id string) error {
	return s.userRepo.Delete(id)
}

func (s *userService) HasAnyUsers() (bool, error) {
	count, err := s.userRepo.Count()
	if err != nil {
		logger.Error("Failed to count users", zap.Error(err))
		return false, err
	}
	return count > 0, nil
}
