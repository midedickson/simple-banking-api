package main

import (
	"errors"
	"log"
	"sync"
)

var ErrInsufficientFunds = errors.New("insufficient funds in account balance")

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// data transfer objects
type createTransactionDTO struct {
	Amount    float64 `json:"amount"`
	AccountID int     `json:"account_id"`
	Direction string  `json:"direction"`
}

type forwardTransactionDTO struct {
	Reference string  `json:"reference"`
	AccountID int     `json:"account_id"`
	Amount    float64 `json:"amount"`
}

// user account details

type UserAccount struct {
	mu        sync.Mutex
	AccountID int     `json:"account_id"`
	Balance   float64 `json:"balance"`
}

func (u *UserAccount) performAccountCredit(amount float64) {
	u.mu.Lock()
	defer u.mu.Unlock()
	log.Printf("Account Balance before credit: %v", u.Balance)
	log.Printf("Crediting: %v", amount)
	u.Balance += amount
	log.Printf("Account Balance after credit: %v", u.Balance)
}

func (u *UserAccount) performAccountDebit(amount float64) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	log.Printf("Account Balance before debit: %v", u.Balance)
	log.Printf("Debiting: %v", amount)
	if u.Balance < amount {
		log.Println("Debit Refused, Insufficient Funds")
		return ErrInsufficientFunds
	}
	u.Balance -= amount
	log.Printf("Account Balance after debit: %v", u.Balance)
	return nil
}
