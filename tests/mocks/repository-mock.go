package mocks

import (
	"github.com/midedickson/simple-banking-app/dto"
	"github.com/midedickson/simple-banking-app/models"
	"github.com/stretchr/testify/mock"
)

// Mock the repository
type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) FindAccountById(userAccountId int) *models.UserAccount {
	args := m.Called(userAccountId)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*models.UserAccount)
}

func (m *MockRepo) CreateTransaction(createTransactionDTO *dto.CreateDBTransactionDTO) (*models.Transaction, error) {
	args := m.Called(createTransactionDTO)
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *MockRepo) UpdateTransactionStatus(transaction *models.Transaction, status string) error {
	args := m.Called(transaction, status)
	return args.Error(0)
}

func (m *MockRepo) FetchTransactionDetailsByReference(reference string) *models.Transaction {
	args := m.Called(reference)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*models.Transaction)
}

func (m *MockRepo) GenerateTransactionReference() string {
	args := m.Called()
	return args.String(0)
}
