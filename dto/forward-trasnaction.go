package dto

type ForwardTransactionDTO struct {
	Reference string  `json:"reference"`
	AccountID int     `json:"account_id"`
	Amount    float64 `json:"amount"`
}
