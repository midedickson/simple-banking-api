package models

import (
	"log"
	"sync"

	"github.com/midedickson/simple-banking-app/constants"
	"github.com/shopspring/decimal"
)

// user account details

type UserAccount struct {
	mu      sync.Mutex
	ID      int             `json:"account_id"`
	Balance decimal.Decimal `json:"balance"`
}

func (u *UserAccount) Credit(amount float64) {
	u.mu.Lock()
	defer u.mu.Unlock()
	log.Printf("Account Balance before credit: %v", u.Balance)
	log.Printf("Crediting: %v", amount)
	u.Balance.Add(decimal.NewFromFloat(amount))
	log.Printf("Account Balance after credit: %v", u.Balance)
}

func (u *UserAccount) Debit(amount float64) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	log.Printf("Account Balance before debit: %v", u.Balance)
	log.Printf("Debiting: %v", amount)
	if u.Balance.LessThan(decimal.NewFromFloat(amount)) {
		log.Println("Debit Refused, Insufficient Funds")
		return constants.ErrInsufficientFunds
	}
	u.Balance.Sub(decimal.NewFromFloat(amount))
	log.Printf("Account Balance after debit: %v", u.Balance)
	return nil
}
