package repository

import (
	"smart-contract-automation/src/internal/model"

	"gorm.io/gorm"
)

type WorkflowRepository interface {
	Create(workflow *model.Workflow) error
	List(page, limit int) ([]model.Contract, error)
	Count() (int64, error)
	GetByID(id string) (*model.Workflow, error)
	Update(workflow *model.Workflow) error
	Delete(id string) error
}

type workflowRepository struct {
	db *gorm.DB
}

func NewWorkflowRepository(db *gorm.DB) WorkflowRepository {
	return &workflowRepository{db: db}
}

func (r *workflowRepository) Create(workflow *model.Workflow) error {
	return r.db.Create(workflow).Error
}

func (r *workflowRepository) List(page, limit int) ([]model.Contract, error) {
	var contracts []model.Contract
	offset := (page - 1) * limit
	err := r.db.Offset(offset).Limit(limit).Find(&contracts).Error
	return contracts, err
}

func (r *workflowRepository) Count() (int64, error) {
	var total int64
	err := r.db.Model(&model.Contract{}).Count(&total).Error
	return total, err
}

func (r *workflowRepository) GetByID(id string) (*model.Workflow, error) {
	var workflow model.Workflow
	if err := r.db.First(&workflow, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &workflow, nil
}

func (r *workflowRepository) Update(workflow *model.Workflow) error {
	return r.db.Save(workflow).Error
}

func (r *workflowRepository) Delete(id string) error {
	return r.db.Delete(&model.Workflow{}, "id = ?", id).Error
}
