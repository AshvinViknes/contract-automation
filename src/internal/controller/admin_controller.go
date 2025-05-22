package controller

import (
	"encoding/json"
	"net/http"
	"smart-contract-automation/src/internal/model"
	"smart-contract-automation/src/internal/service"
	"smart-contract-automation/src/pkg/common"
	"smart-contract-automation/src/pkg/logger"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type AdminController struct {
	userService service.UserService
}

type CreateAdminRequest struct {
	Email     string     `json:"email" binding:"required,email"`
	Password  string     `json:"password" binding:"required,min=8"`
	Role      model.Role `json:"role" binding:"required"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
}

func NewAdminController(userService service.UserService) *AdminController {
	return &AdminController{
		userService: userService,
	}
}

// CreateAdmin creates a new admin or super admin user
func (c *AdminController) CreateAdmin(w http.ResponseWriter, r *http.Request) {
	var req CreateAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate role
	if req.Role != model.RoleAdmin && req.Role != model.RoleSuperAdmin {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid role. Must be 'admin' or 'super_admin'")
		return
	}

	user := &model.User{
		Email:     req.Email,
		Role:      req.Role,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	if err := c.userService.CreateUser(user, req.Password); err != nil {
		logger.Error("Failed to create admin user", zap.Error(err))
		common.RespondWithError(w, http.StatusInternalServerError, "Failed to create admin user")
		return
	}

	// Don't return the password hash in the response
	user.Password = ""
	common.RespondWithSuccess(w, http.StatusCreated, user)
}

func (c *AdminController) GetAdmin(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid UserID")
		return
	}

	user, err := c.userService.GetUser(id)
	if err != nil {
		common.RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	common.RespondWithSuccess(w, http.StatusOK, user)
}

// ListAdmins returns a list of all admin users
func (c *AdminController) ListAdmins(w http.ResponseWriter, r *http.Request) {
	admins, err := c.userService.ListAdmins()
	if err != nil {
		logger.Error("Failed to list admin users", zap.Error(err))
		common.RespondWithError(w, http.StatusInternalServerError, "Failed to list admin users")
		return
	}

	common.RespondWithSuccess(w, http.StatusOK, admins)
}

// DeleteAdmin deletes an admin user
func (c *AdminController) DeleteAdmin(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		common.RespondWithError(w, http.StatusBadRequest, "UserID is required")
		return
	}

	user, err := c.userService.GetUser(id)
	if err != nil {
		common.RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	if err := c.userService.DeleteUser(user.ID); err != nil {
		logger.Error("Failed to delete admin user", zap.Error(err), zap.String("id", user.ID))
		common.RespondWithError(w, http.StatusInternalServerError, "Failed to delete admin user")
		return
	}

	common.RespondWithSuccess(w, http.StatusOK, map[string]string{"message": "Admin user deleted successfully"})
}
