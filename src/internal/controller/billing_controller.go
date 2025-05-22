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

type BillingController struct {
	service service.BillingService
}

func NewBillingController(service service.BillingService) *BillingController {
	return &BillingController{service: service}
}

// POST /api/v1/billing
func (bc *BillingController) CreateBill(w http.ResponseWriter, r *http.Request) {
	var bill model.Billing
	if err := json.NewDecoder(r.Body).Decode(&bill); err != nil {
		logger.Error("Failed to decode bill request", zap.Error(err))
		common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := bc.service.CreateBill(&bill); err != nil {
		logger.Error("Failed to create bill", zap.Error(err))
		common.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	logger.Info("Billing created successfully", zap.String("id", bill.ID))
	common.RespondWithSuccess(w, http.StatusCreated, bill)
}

// GET /api/v1/billing/{id}
func (bc *BillingController) GetBill(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid billing ID")
		return
	}

	bill, err := bc.service.GetBill(id)
	if err != nil {
		common.RespondWithError(w, http.StatusNotFound, "Billing not found")
		return
	}

	common.RespondWithSuccess(w, http.StatusOK, bill)
}

// PUT /api/v1/billing/{id}
func (bc *BillingController) UpdateBill(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid billing ID")
		return
	}

	var bill model.Billing
	if err := json.NewDecoder(r.Body).Decode(&bill); err != nil {
		logger.Error("Failed to decode bill request", zap.Error(err))
		common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	bill.ID = id
	if err := bc.service.UpdateBill(&bill); err != nil {
		logger.Error("Failed to update billing", zap.Error(err))
		common.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	common.RespondWithSuccess(w, http.StatusOK, bill)
}

// DELETE /api/v1/billing/{id}
func (bc *BillingController) DeleteBill(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid billing ID")
		return
	}

	if err := bc.service.DeleteBill(id); err != nil {
		logger.Error("Failed to delete billing", zap.Error(err))
		common.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	common.RespondWithSuccess(w, http.StatusOK, map[string]string{"message": "Billing deleted successfully"})
}

// GET /api/v1/billing
func (bc *BillingController) ListBills(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page == 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 10
	}

	bills, total, err := bc.service.ListBills(page, limit)
	if err != nil {
		common.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	common.RespondWithPaginatedSuccess(w, http.StatusOK, bills, page, limit, total)
}

func (bc *BillingController) ListBillsByUserID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Missing UserID")
		return
	}

	bills, err := bc.service.ListBillsByUserID(userID)
	if err != nil {
		logger.Error("Failed to list bills by userID", zap.Error(err))
		common.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	common.RespondWithSuccess(w, http.StatusOK, bills)
}
