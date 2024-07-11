package dto

// data transfer object for creating transaction
type CreateTransactionDTO struct {
	Amount    float64 `json:"amount"`
	AccountID int     `json:"account_id"`
	Direction string  `json:"direction"`
}
