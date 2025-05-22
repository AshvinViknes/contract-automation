package model

import "time"

type Billing struct {
	ID            string    `json:"id" gorm:"primaryKey;type:varchar(21)"`
	UserID        string    `json:"user_id"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	Status        string    `json:"status"`
	BillingDate   time.Time `json:"billing_date"`
	DueDate       time.Time `json:"due_date"`
	PaymentMethod string    `json:"payment_method"`
	TransactionID string    `json:"transaction_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
