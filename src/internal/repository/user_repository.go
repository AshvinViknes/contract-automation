package repository

import (
	"smart-contract-automation/src/internal/model"
	"smart-contract-automation/src/pkg/utils"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetByEmail(email string) (*model.User, error)
	Create(user *model.User) error
	GetByID(id string) (*model.User, error)
	FindByRoles(roles []model.Role) ([]model.User, error)
	Delete(id string) error
	Count() (int64, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(user *model.User) error {
	// Generate ID before creating user
	id := utils.GenerateID()
	user.ID = id
	return r.db.Create(user).Error
}

func (r *userRepository) FindByRoles(roles []model.Role) ([]model.User, error) {
	var users []model.User
	err := r.db.Where("role IN ?", roles).Find(&users).Error
	return users, err
}

func (r *userRepository) GetByID(id string) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.User{}).Error
}

func (r *userRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.User{}).Count(&count).Error
	return count, err
}
