package repository

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/midedickson/simple-banking-app/dto"
	"github.com/midedickson/simple-banking-app/models"
	"gorm.io/gorm"
)

type StorageRepository struct {
	DB *gorm.DB
}

func NewStorageRepository(DB *gorm.DB) *StorageRepository {
	return &StorageRepository{DB: DB}
}

func (r *StorageRepository) GenerateTransactionReference() string {
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

func (r *StorageRepository) FindAccountById(userAccountId int) *models.UserAccount {
	// todo: implement logic to find the account based on UserAccountId
	// iterate through the users slice and find the account with matching UserAccountId
	for _, user := range Users {
		if user.ID == userAccountId {
			return user
		}
	}

	// if no account found, return nil
	return nil
}

func (r *StorageRepository) CreateTransaction(createTransactionDTO *dto.CreateDBTransactionDTO) (*models.Transaction, error) {
	transaction := models.Transaction{
		AccountID: createTransactionDTO.AccountID,
		Reference: r.GenerateTransactionReference(),
		Amount:    createTransactionDTO.Amount,
		Status:    "pending",
		Direction: createTransactionDTO.Direction,
	}

	return &transaction, r.DB.Create(&transaction).Error
}

func (r *StorageRepository) UpdateTransactionStatus(transaction *models.Transaction, status string) error {
	transaction.Status = status
	return r.DB.Save(&transaction).Error
}

func (r *StorageRepository) FetchTransactionDetailsByReference(reference string) *models.Transaction {
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
