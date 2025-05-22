package service

import (
	"errors"
	"fmt"
	"smart-contract-automation/src/internal/model"
	"smart-contract-automation/src/pkg/logger"
	"time"

	"smart-contract-automation/src/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
)

type AuthService interface {
	Login(email, password string) (*model.TokenResponse, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
	HashPassword(password string) (string, error)
	ComparePasswords(hashedPassword, password string) error
}

type authService struct {
	userRepo repository.UserRepository
	jwtKey   []byte
}

type Claims struct {
	UserID string     `json:"user_id"`
	Email  string     `json:"email"`
	Role   model.Role `json:"role"`
	jwt.RegisteredClaims
}

func NewAuthService(userRepo repository.UserRepository, jwtKey string) AuthService {
	return &authService{
		userRepo: userRepo,
		jwtKey:   []byte(jwtKey),
	}
}

func (s *authService) Login(email, password string) (*model.TokenResponse, error) {
	logger.Debug("Attempting login", zap.String("email", email))

	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		logger.Error("Failed to get user", zap.Error(err), zap.String("email", email))
		return nil, ErrInvalidCredentials
	}

	if err := s.ComparePasswords(user.Password, password); err != nil {
		return nil, ErrInvalidCredentials
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtKey)
	if err != nil {
		logger.Error("Failed to sign token", zap.Error(err))
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &model.TokenResponse{
		AccessToken: tokenString,
		TokenType:   "Bearer",
		ExpiresAt:   expirationTime,
	}, nil
}

func (s *authService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtKey, nil
	})
}

func (s *authService) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (s *authService) ComparePasswords(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
