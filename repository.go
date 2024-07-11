package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

type repository struct {
	DB *gorm.DB
}

func newRepository() *repository {
	return &repository{DB: DB}
}

func (r *repository) generateTransactionReference() string {
	// Generate unique transaction reference
	// Get the current time as Unix timestamp
	timestamp := time.Now().UnixNano()

	// Generate a random number
	randomNum := rand.Int63()

	// Create a transaction ID using prefix, timestamp, and random number
	transactionID := fmt.Sprintf("%s-%d-%d", "TRX", timestamp, randomNum)

	// check if the transaction exists before returning final value
	if existingTransaction := r.fetchTransactionDetailsByReference(transactionID); existingTransaction != nil {
		return r.generateTransactionReference()
	}
	return transactionID
}

func (r *repository) findAccountById(userAccountId int) *UserAccount {
	// todo: implement logic to find the account based on UserAccountId
	// iterate through the users slice and find the account with matching UserAccountId
	for _, user := range users {
		if user.AccountID == userAccountId {
			return user
		}
	}

	// if no account found, return nil
	return nil
}

func (r *repository) newTransaction(createTransactionDTO *createTransactionDTO) (*Transaction, error) {
	transaction := Transaction{
		AccountID: createTransactionDTO.AccountID,
		Reference: r.generateTransactionReference(),
		Amount:    createTransactionDTO.Amount,
		Status:    "pending",
		Direction: createTransactionDTO.Direction,
	}

	return &transaction, r.DB.Create(&transaction).Error
}

func (r *repository) updateTransactionStatus(transaction *Transaction, status string) {
	transaction.Status = status
	r.DB.Save(&transaction)
}

func (r *repository) fetchTransactionDetailsByReference(reference string) *Transaction {
	var transaction Transaction
	result := r.DB.Where("reference =?", reference).First(&transaction)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil // No record found, return nil
		}
		log.Println("Error fetching transaction:", result.Error)
	}
	if transaction.Reference != "" {
		return &transaction // Transaction found, return reference
	}
	return nil // No transaction found, return nil
}
