package controller

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"smart-contract-automation/src/internal/model"
	"smart-contract-automation/src/internal/service"
	"smart-contract-automation/src/pkg/common"
	"smart-contract-automation/src/pkg/config"
	"smart-contract-automation/src/pkg/logger"

	"go.uber.org/zap"
)

type SetupController struct {
	userService service.UserService
	setupToken  string
}

type setupRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

func (r *setupRequest) validate() error {
	if r.Email == "" || r.Password == "" || r.FirstName == "" || r.LastName == "" {
		return common.NewValidationError("all fields are required")
	}
	if len(r.Password) < 8 {
		return common.NewValidationError("password must be at least 8 characters")
	}
	if _, err := mail.ParseAddress(r.Email); err != nil {
		return common.NewValidationError("invalid email format")
	}
	return nil
}

func NewSetupController(userService service.UserService, cfg *config.Config) *SetupController {
	if cfg.SetupToken == "" {
		logger.Error("Setup token is not configured")
	}
	return &SetupController{
		userService: userService,
		setupToken:  cfg.SetupToken,
	}
}

func (c *SetupController) CreateFirstSuperAdmin(w http.ResponseWriter, r *http.Request) {
	// Check if setup is enabled
	if c.setupToken == "" {
		logger.Error("Setup endpoint is disabled")
		common.RespondWithError(w, http.StatusForbidden, "Setup endpoint is disabled")
		return
	}

	// Verify setup token
	setupToken := r.Header.Get("X-Setup-Token")
	if setupToken == "" {
		logger.Error("Missing setup token")
		common.RespondWithError(w, http.StatusUnauthorized, "Setup token is required")
		return
	}

	if setupToken != c.setupToken {
		logger.Error("Invalid setup token attempt")
		common.RespondWithError(w, http.StatusUnauthorized, "Invalid setup token")
		return
	}

	// Parse and validate request
	var req setupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Failed to decode request body", zap.Error(err))
		common.RespondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}
	defer r.Body.Close()

	if err := req.validate(); err != nil {
		logger.Error("Invalid request data", zap.Error(err))
		common.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Check if any users exist
	exists, err := c.userService.HasAnyUsers()
	if err != nil {
		logger.Error("Failed to check existing users", zap.Error(err))
		common.RespondWithError(w, http.StatusInternalServerError, "Failed to verify system state")
		return
	}

	if exists {
		logger.Error("Attempt to create first super admin when users already exist")
		common.RespondWithError(w, http.StatusConflict, "Super admin already exists")
		return
	}

	// Create super admin user
	user := &model.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      model.RoleSuperAdmin,
	}

	if err := c.userService.CreateUser(user, req.Password); err != nil {
		logger.Error("Failed to create first super admin", zap.Error(err))
		common.RespondWithError(w, http.StatusInternalServerError, "Failed to create super admin")
		return
	}

	// Don't return the password hash
	user.Password = ""
	logger.Info("Super admin created successfully", zap.String("email", user.Email))
	common.RespondWithSuccess(w, http.StatusCreated, user)
}
