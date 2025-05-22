package service

import (
	"smart-contract-automation/src/internal/model"
	"smart-contract-automation/src/internal/repository"
	"smart-contract-automation/src/pkg/logger"

	"go.uber.org/zap"
)

type WorkflowService interface {
	CreateWorkflow(workflow *model.Workflow) error
	ListWorkflows(page, limit int) ([]model.Contract, int64, error)
	GetWorkflow(id string) (*model.Workflow, error)
	ApproveWorkflow(id string) error
	RejectWorkflow(id string) error
}

type workflowService struct {
	repo repository.WorkflowRepository
}

func NewWorkflowService(repo repository.WorkflowRepository) WorkflowService {
	return &workflowService{repo: repo}
}

func (s *workflowService) CreateWorkflow(workflow *model.Workflow) error {
	// Set initial status
	workflow.ID = "wf_" + workflow.ContractID
	workflow.Status = "Pending"
	return s.repo.Create(workflow)
}

func (s *workflowService) GetWorkflow(id string) (*model.Workflow, error) {
	return s.repo.GetByID(id)
}

func (s *workflowService) ListWorkflows(page, limit int) ([]model.Contract, int64, error) {
	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}

	contracts, err := s.repo.List(page, limit)
	return contracts, total, err
}

func (s *workflowService) ApproveWorkflow(id string) error {
	workflow, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	workflow.Status = "Approved"
	if err := s.repo.Update(workflow); err != nil {
		logger.Error("Failed to update workflow", zap.Error(err))
		return err
	}
	return nil
}

func (s *workflowService) RejectWorkflow(id string) error {
	workflow, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	workflow.Status = "Rejected"
	return s.repo.Update(workflow)
}
