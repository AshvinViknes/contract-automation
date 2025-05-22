package controller

import (
	"encoding/json"
	"net/http"
	"smart-contract-automation/src/internal/model"
	"smart-contract-automation/src/internal/service"
	"smart-contract-automation/src/pkg/common"
	"smart-contract-automation/src/pkg/logger"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type WorkflowController struct {
	service service.WorkflowService
}

func NewWorkflowController(service service.WorkflowService) *WorkflowController {
	return &WorkflowController{service: service}
}

// POST /api/v1/workflows
func (wc *WorkflowController) CreateWorkflow(w http.ResponseWriter, r *http.Request) {
	var workflow model.Workflow
	if err := json.NewDecoder(r.Body).Decode(&workflow); err != nil {
		logger.Error("Failed to decode workflow request", zap.Error(err))
		common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := wc.service.CreateWorkflow(&workflow); err != nil {
		logger.Error("Failed to create workflow", zap.Error(err))
		common.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	logger.Info("Workflow created successfully", zap.String("id", workflow.ID))
	common.RespondWithSuccess(w, http.StatusCreated, workflow)
}

// GET /api/v1/workflows/{id}
func (wc *WorkflowController) GetWorkflow(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid workflow ID")
		return
	}

	workflow, err := wc.service.GetWorkflow(id)
	if err != nil {
		common.RespondWithError(w, http.StatusNotFound, "Workflow not found")
		return
	}

	common.RespondWithSuccess(w, http.StatusOK, workflow)
}

// GET /api/v1/workflows
func (wc *WorkflowController) ListWorkflows(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page == 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 10
	}

	contracts, total, err := wc.service.ListWorkflows(page, limit)
	if err != nil {
		common.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	common.RespondWithPaginatedSuccess(w, http.StatusOK, contracts, page, limit, total)
}

// PUT /api/v1/workflows/{id}/approve
func (wc *WorkflowController) ApproveWorkflow(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid workflow ID")
		return
	}

	if err := wc.service.ApproveWorkflow(id); err != nil {
		logger.Error("Failed to approve workflow", zap.Error(err))
		common.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	common.RespondWithSuccess(w, http.StatusOK, map[string]string{"message": "Workflow approved successfully"})
}

// PUT /api/v1/workflows/{id}/reject
func (wc *WorkflowController) RejectWorkflow(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid workflow ID")
		return
	}

	if err := wc.service.RejectWorkflow(id); err != nil {
		logger.Error("Failed to reject workflow", zap.Error(err))
		common.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	common.RespondWithSuccess(w, http.StatusOK, map[string]string{"message": "Workflow rejected successfully"})
}
