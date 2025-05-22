package model

import (
	"time"

	"gorm.io/gorm"
)

type Contract struct {
	ID             string    `json:"id" gorm:"primaryKey;type:varchar(21)"`
	Name           string    `json:"name"`
	Content        string    `json:"content"`
	Status         string    `json:"status"`
	Conditions     string    `json:"conditions"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	CreatedBy      string    `json:"created_by" gorm:"type:varchar(21)"`
	LastModifiedBy string    `json:"last_modified_by" gorm:"type:varchar(21)"`
}

type ContractTemplate struct {
	gorm.Model
	Name     string `json:"name"`
	Content  string `json:"content"`
	Category string `json:"category"`
	IsActive bool   `json:"is_active"`
}

type StatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

type ContractTrend struct {
	Date   string `json:"date"`
	Status string `json:"status"`
	Count  int64  `json:"count"`
}
