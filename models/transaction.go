package models

import (
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	AccountID int     `gorm:"account_id" json:"account_id"`
	Reference string  `gorm:"reference" json:"reference"`
	Amount    float64 `gorm:"amount" json:"amount"`
	Direction string  `gorm:"direction" json:"direction"`
	Status    string  `gorm:"status" json:"status"`
}
