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

type ContractController struct {
	service service.ContractService
}

func NewContractController(service service.ContractService) *ContractController {
	return &ContractController{service: service}
}

func (c *ContractController) CreateContract(w http.ResponseWriter, r *http.Request) {
	var contract model.Contract
	if err := json.NewDecoder(r.Body).Decode(&contract); err != nil {
		logger.Error("Failed to decode contract request", zap.Error(err))
		common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := c.service.CreateContract(&contract); err != nil {
		logger.Error("Failed to create contract", zap.Error(err))
		common.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	logger.Info("Contract created successfully", zap.String("id", contract.ID))
	common.RespondWithSuccess(w, http.StatusCreated, contract)
}

func (c *ContractController) GetContract(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid contract ID")
		return
	}

	contract, err := c.service.GetContract(id)
	if err != nil {
		common.RespondWithError(w, http.StatusNotFound, "Contract not found")
		return
	}

	common.RespondWithSuccess(w, http.StatusOK, contract)
}

func (c *ContractController) ListContracts(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page == 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 10
	}

	contracts, total, err := c.service.ListContracts(page, limit)
	if err != nil {
		common.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	common.RespondWithPaginatedSuccess(w, http.StatusOK, contracts, page, limit, total)
}

func (c *ContractController) UpdateContract(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid contract ID")
		return
	}

	var contract model.Contract
	if err := json.NewDecoder(r.Body).Decode(&contract); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	contract.ID = id
	if err := c.service.UpdateContract(&contract); err != nil {
		common.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	common.RespondWithSuccess(w, http.StatusOK, contract)
}

func (c *ContractController) DeleteContract(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid contract ID")
		return
	}

	if err := c.service.DeleteContract(id); err != nil {
		common.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	common.RespondWithSuccess(w, http.StatusOK, map[string]string{"message": "Contract deleted successfully"})
}

func (c *ContractController) DeployContract(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid contract ID")
		return
	}

	// Call service to deploy the contract
	txHash, err := c.service.DeployContract(id)
	if err != nil {
		logger.Error("Failed to deploy contract", zap.Error(err))
		common.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	common.RespondWithSuccess(w, http.StatusOK, map[string]string{
		"message": "Contract deployed successfully",
		"txHash":  txHash,
	})
}

func (c *ContractController) TerminateContract(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid contract ID")
		return
	}

	if err := c.service.TerminateContract(id); err != nil {
		logger.Error("Failed to terminate contract", zap.Error(err))
		common.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	common.RespondWithSuccess(w, http.StatusOK, map[string]string{
		"message": "Contract terminated successfully",
	})
}

func (c *ContractController) GetContractStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid contract ID")
		return
	}

	status, err := c.service.GetContractStatus(id)
	if err != nil {
		common.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	common.RespondWithSuccess(w, http.StatusOK, map[string]string{
		"status": status,
	})
}

// GET /api/v1/analytics/status/count
func (c *ContractController) GetContractStatusCounts(w http.ResponseWriter, r *http.Request) {
	counts, err := c.service.GetContractStatusCounts()
	if err != nil {
		logger.Error("Failed to get contract status counts", zap.Error(err))
		common.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	common.RespondWithSuccess(w, http.StatusOK, counts)
}

// GET /api/v1/analytics/contracts?status=Pending
func (c *ContractController) GetContractsByStatus(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	if status == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Missing status parameter")
		return
	}
	contracts, err := c.service.GetContractsByStatus(status)
	if err != nil {
		logger.Error("Failed to get contracts by status", zap.Error(err))
		common.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	common.RespondWithSuccess(w, http.StatusOK, contracts)
}

// GET /api/v1/analytics/trend
func (c *ContractController) GetContractTrend(w http.ResponseWriter, r *http.Request) {
	granularity := r.URL.Query().Get("granularity")
	trends, err := c.service.GetContractTrend(granularity)
	if err != nil {
		logger.Error("Failed to get contract trends", zap.Error(err))
		common.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	common.RespondWithSuccess(w, http.StatusOK, trends)
}
