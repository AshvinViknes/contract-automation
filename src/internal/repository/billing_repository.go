package repository

import (
	"smart-contract-automation/src/internal/model"
	"smart-contract-automation/src/pkg/utils"

	"gorm.io/gorm"
)

type BillingRepository interface {
	Create(bill *model.Billing) error
	GetByID(id string) (*model.Billing, error)
	Update(bill *model.Billing) error
	Delete(id string) error
	List(page, limit int) ([]model.Billing, error)
	Count() (int64, error)
	ListByUserID(userID string) ([]model.Billing, error)
}

type billingRepository struct {
	db *gorm.DB
}

func NewBillingRepository(db *gorm.DB) BillingRepository {
	return &billingRepository{db: db}
}

func (r *billingRepository) Create(bill *model.Billing) error {
	id := utils.GenerateID()
	bill.ID = id
	return r.db.Create(bill).Error
}

func (r *billingRepository) GetByID(id string) (*model.Billing, error) {
	var bill model.Billing
	err := r.db.First(&bill, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &bill, nil
}

func (r *billingRepository) Update(bill *model.Billing) error {
	return r.db.Save(bill).Error
}

func (r *billingRepository) Delete(id string) error {
	return r.db.Delete(&model.Billing{}, "id = ?", id).Error
}

func (r *billingRepository) List(page, limit int) ([]model.Billing, error) {
	var bills []model.Billing
	offset := (page - 1) * limit
	err := r.db.Offset(offset).Limit(limit).Find(&bills).Error
	return bills, err
}

func (r *billingRepository) Count() (int64, error) {
	var total int64
	err := r.db.Model(&model.Billing{}).Count(&total).Error
	return total, err
}

func (r *billingRepository) ListByUserID(userID string) ([]model.Billing, error) {
	var bills []model.Billing
	err := r.db.Where("user_id = ?", userID).Find(&bills).Error
	return bills, err
}
