package service

import (
	"smart-contract-automation/src/internal/model"
	"smart-contract-automation/src/internal/repository"
	"smart-contract-automation/src/pkg/logger"
	"smart-contract-automation/src/pkg/utils"

	"go.uber.org/zap"
)

type ContractService interface {
	CreateContract(contract *model.Contract) error
	GetContract(id string) (*model.Contract, error)
	UpdateContract(contract *model.Contract) error
	DeleteContract(id string) error
	ListContracts(page, limit int) ([]model.Contract, int64, error)
	DeployContract(id string) (string, error)
	TerminateContract(id string) error
	GetContractStatus(id string) (string, error)
	GetContractStatusCounts() (map[string]int64, error)
	GetContractsByStatus(status string) ([]model.Contract, error)
	GetContractTrend(granularity string) ([]model.ContractTrend, error)
}

type contractService struct {
	repo repository.ContractRepository
}

func NewContractService(repo repository.ContractRepository) ContractService {
	return &contractService{repo: repo}
}

func (s *contractService) CreateContract(contract *model.Contract) error {
	// Generate ID for the contract
	id := utils.GenerateID()
	contract.ID = id
	contract.Status = "Pending"

	if err := s.repo.Create(contract); err != nil {
		logger.Error("Failed to create contract in repository", zap.Error(err))
		return err
	}

	return nil
}

func (s *contractService) GetContract(id string) (*model.Contract, error) {
	return s.repo.GetByID(id)
}

func (s *contractService) UpdateContract(contract *model.Contract) error {
	return s.repo.Update(contract)
}

func (s *contractService) DeleteContract(id string) error {
	return s.repo.Delete(id)
}

func (s *contractService) ListContracts(page, limit int) ([]model.Contract, int64, error) {
	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}

	contracts, err := s.repo.List(page, limit)
	return contracts, total, err
}

func (s *contractService) DeployContract(id string) (string, error) {
	contract, err := s.repo.GetByID(id)
	if err != nil {
		return "", err
	}

	// Mock deployment logic
	txHash := utils.GenerateID() // Replace with actual deployment logic

	// Update contract status after deployment
	contract.Status = "Deployed"
	if err := s.repo.Update(contract); err != nil {
		return "", err
	}

	return txHash, nil
}

func (s *contractService) TerminateContract(id string) error {
	contract, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	contract.Status = "Terminated"
	return s.repo.Update(contract)
}

func (s *contractService) GetContractStatus(id string) (string, error) {
	contract, err := s.repo.GetByID(id)
	if err != nil {
		return "", err
	}
	return contract.Status, nil
}

func (s *contractService) GetContractStatusCounts() (map[string]int64, error) {
	counts, err := s.repo.CountByStatus()
	if err != nil {
		return nil, err
	}
	result := map[string]int64{}
	for _, c := range counts {
		result[c.Status] = c.Count
	}
	return result, nil
}

func (s *contractService) GetContractsByStatus(status string) ([]model.Contract, error) {
	return s.repo.ListByStatus(status)
}

func (s *contractService) GetContractTrend(granularity string) ([]model.ContractTrend, error) {
	return s.repo.GetContractTrend(granularity)
}
