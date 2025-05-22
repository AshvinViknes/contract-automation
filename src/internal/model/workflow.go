package model

import "time"

type Workflow struct {
	ID         string    `json:"id" gorm:"primaryKey;type:varchar(21)"`
	ContractID string    `json:"contract_id"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
