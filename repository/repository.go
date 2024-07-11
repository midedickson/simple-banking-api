package repository

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/midedickson/simple-banking-app/config"
	"github.com/midedickson/simple-banking-app/dto"
	"github.com/midedickson/simple-banking-app/models"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository() *Repository {
	return &Repository{DB: config.DB}
}

func (r *Repository) GenerateTransactionReference() string {
	// Generate unique transaction reference
	// Get the current time as Unix timestamp
	timestamp := time.Now().UnixNano()

	// Generate a random number
	randomNum := rand.Int63()

	// Create a transaction ID using prefix, timestamp, and random number
	transactionID := fmt.Sprintf("%s-%d-%d", "TRX", timestamp, randomNum)

	// check if the transaction exists before returning final value
	if existingTransaction := r.FetchTransactionDetailsByReference(transactionID); existingTransaction != nil {
		return r.GenerateTransactionReference()
	}
	return transactionID
}

func (r *Repository) FindAccountById(userAccountId int) *models.UserAccount {
	// todo: implement logic to find the account based on UserAccountId
	// iterate through the users slice and find the account with matching UserAccountId
	for _, user := range Users {
		if user.AccountID == userAccountId {
			return user
		}
	}

	// if no account found, return nil
	return nil
}

func (r *Repository) NewTransaction(createTransactionDTO *dto.CreateTransactionDTO) (*models.Transaction, error) {
	transaction := models.Transaction{
		AccountID: createTransactionDTO.AccountID,
		Reference: r.GenerateTransactionReference(),
		Amount:    createTransactionDTO.Amount,
		Status:    "pending",
		Direction: createTransactionDTO.Direction,
	}

	return &transaction, r.DB.Create(&transaction).Error
}

func (r *Repository) UpdateTransactionStatus(transaction *models.Transaction, status string) {
	transaction.Status = status
	r.DB.Save(&transaction)
}

func (r *Repository) FetchTransactionDetailsByReference(reference string) *models.Transaction {
	var transaction models.Transaction
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
