package service

import (
	"smart-contract-automation/src/pkg/logger"
	"smart-contract-automation/src/pkg/utils"

	"smart-contract-automation/src/internal/model"
	"smart-contract-automation/src/internal/repository"

	"go.uber.org/zap"
)

type BillingService interface {
	CreateBill(bill *model.Billing) error
	GetBill(id string) (*model.Billing, error)
	UpdateBill(bill *model.Billing) error
	DeleteBill(id string) error
	ListBills(page, limit int) ([]model.Billing, int64, error)
	ListBillsByUserID(userID string) ([]model.Billing, error)
}

type billingService struct {
	repo repository.BillingRepository
}

func NewBillingService(repo repository.BillingRepository) BillingService {
	return &billingService{repo: repo}
}

func (s *billingService) CreateBill(bill *model.Billing) error {
	// Generate ID for the billing
	id := utils.GenerateID()
	bill.ID = id
	bill.Status = "Pending"

	if err := s.repo.Create(bill); err != nil {
		logger.Error("Failed to create billing in repository", zap.Error(err))
		return err
	}

	return nil
}

func (s *billingService) GetBill(id string) (*model.Billing, error) {
	return s.repo.GetByID(id)
}

func (s *billingService) UpdateBill(bill *model.Billing) error {
	return s.repo.Update(bill)
}

func (s *billingService) DeleteBill(id string) error {
	return s.repo.Delete(id)
}

func (s *billingService) ListBills(page, limit int) ([]model.Billing, int64, error) {
	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}

	bills, err := s.repo.List(page, limit)
	return bills, total, err
}

func (s *billingService) ListBillsByUserID(userID string) ([]model.Billing, error) {
	return s.repo.ListByUserID(userID)
}
