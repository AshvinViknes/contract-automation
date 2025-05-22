package repository

import (
	"smart-contract-automation/src/internal/model"
	"smart-contract-automation/src/pkg/utils"

	"gorm.io/gorm"
)

type ContractRepository interface {
	Create(contract *model.Contract) error
	GetByID(id string) (*model.Contract, error)
	Update(contract *model.Contract) error
	Delete(id string) error
	List(page, limit int) ([]model.Contract, error)
	Count() (int64, error)
	CountByStatus() ([]model.StatusCount, error)
	ListByStatus(status string) ([]model.Contract, error)
	GetContractTrend(granularity string) ([]model.ContractTrend, error)
}

type contractRepository struct {
	db *gorm.DB
}

func NewContractRepository(db *gorm.DB) ContractRepository {
	return &contractRepository{db: db}
}

func (r *contractRepository) Create(contract *model.Contract) error {
	// Generate ID before creating contract
	id := utils.GenerateID()
	contract.ID = id
	return r.db.Create(contract).Error
}

func (r *contractRepository) GetByID(id string) (*model.Contract, error) {
	var contract model.Contract
	err := r.db.First(&contract, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &contract, nil
}

func (r *contractRepository) Update(contract *model.Contract) error {
	return r.db.Save(contract).Error
}

func (r *contractRepository) Delete(id string) error {
	return r.db.Delete(&model.Contract{}, "id = ?", id).Error
}

func (r *contractRepository) List(page, limit int) ([]model.Contract, error) {
	var contracts []model.Contract
	offset := (page - 1) * limit
	err := r.db.Offset(offset).Limit(limit).Find(&contracts).Error
	return contracts, err
}

func (r *contractRepository) Count() (int64, error) {
	var total int64
	err := r.db.Model(&model.Contract{}).Count(&total).Error
	return total, err
}

// CountByStatus uses GROUP BY to get counts for each status.
func (r *contractRepository) CountByStatus() ([]model.StatusCount, error) {
	var counts []model.StatusCount
	err := r.db.Model(&model.Contract{}).
		Select("status, count(*) as count").
		Where("status IN (?)", []string{"Pending", "Deployed", "Terminated"}).
		Group("status").
		Scan(&counts).Error
	return counts, err
}

// ListByStatus returns contracts filtered by the given status.
func (r *contractRepository) ListByStatus(status string) ([]model.Contract, error) {
	var contracts []model.Contract
	err := r.db.Where("status = ?", status).Find(&contracts).Error
	return contracts, err
}

func (r *contractRepository) GetContractTrend(granularity string) ([]model.ContractTrend, error) {
	var trends []model.ContractTrend

	// Define date grouping format based on granularity.
	// Adjust these expressions based on your database dialect.
	var dateFormat string
	switch granularity {
	case "weekly":
		dateFormat = "DATE_FORMAT(created_at, '%Y-%u')" // Year and week number
	case "monthly":
		dateFormat = "DATE_FORMAT(created_at, '%Y-%m')"
	default: // "daily" or unspecified
		dateFormat = "DATE(created_at)"
	}

	err := r.db.Model(&model.Contract{}).
		Select(dateFormat+" as date, status, count(*) as count").
		Where("status IN (?)", []string{"Pending", "Deployed", "Terminated"}).
		Group(dateFormat + ", status").
		Order("date asc").
		Scan(&trends).Error
	return trends, err
}
