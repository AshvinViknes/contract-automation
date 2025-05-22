package controller

import (
	"encoding/json"
	"net/http"
	"smart-contract-automation/src/internal/model"
	"smart-contract-automation/src/internal/service"
	"smart-contract-automation/src/pkg/common"
)

type AuthController struct {
	AuthService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{
		AuthService: authService,
	}
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var loginReq model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	token, err := c.AuthService.Login(loginReq.Email, loginReq.Password)
	if err != nil {
		common.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	common.RespondWithSuccess(w, http.StatusOK, token)
}
