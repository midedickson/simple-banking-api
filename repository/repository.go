package repository

import (
	"github.com/midedickson/simple-banking-app/dto"
	"github.com/midedickson/simple-banking-app/models"
)

type Repository interface {
	GenerateTransactionReference() string
	CreateTransaction(createTransactionDTO *dto.CreateDBTransactionDTO) (*models.Transaction, error)
	UpdateTransactionStatus(transaction *models.Transaction, status string) error
	FetchTransactionDetailsByReference(reference string) *models.Transaction
	FindAccountById(userAccountId int) *models.UserAccount
}

func NewRepository(repository *Repository) *Repository {
	return repository
}
