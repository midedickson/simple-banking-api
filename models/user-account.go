package models

import (
	"log"
	"sync"

	"github.com/midedickson/simple-banking-app/constants"
)

// user account details

type UserAccount struct {
	mu        sync.Mutex
	AccountID int     `json:"account_id"`
	Balance   float64 `json:"balance"`
}

func (u *UserAccount) Credit(amount float64) {
	u.mu.Lock()
	defer u.mu.Unlock()
	log.Printf("Account Balance before credit: %v", u.Balance)
	log.Printf("Crediting: %v", amount)
	u.Balance += amount
	log.Printf("Account Balance after credit: %v", u.Balance)
}

func (u *UserAccount) Debit(amount float64) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	log.Printf("Account Balance before debit: %v", u.Balance)
	log.Printf("Debiting: %v", amount)
	if u.Balance < amount {
		log.Println("Debit Refused, Insufficient Funds")
		return constants.ErrInsufficientFunds
	}
	u.Balance -= amount
	log.Printf("Account Balance after debit: %v", u.Balance)
	return nil
}
