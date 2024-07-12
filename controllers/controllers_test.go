package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/midedickson/simple-banking-app/dto"
	"github.com/midedickson/simple-banking-app/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
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

func (m *MockRepo) CreateTransaction(createTransactionDTO *dto.CreateTransactionDTO) (*models.Transaction, error) {
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

// Mock the external service
type MockExternal struct {
	mock.Mock
}

func (m *MockExternal) ForwardTransactionToThirdParty(transaction *models.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

func TestCreateTransaction(t *testing.T) {
	mockRepo := new(MockRepo)
	mockExternal := new(MockExternal)
	ctrl := &Controller{repo: mockRepo, external: mockExternal}

	handler := http.HandlerFunc(ctrl.CreateTransaction)

	t.Run("successful transaction", func(t *testing.T) {
		transactionDTO := dto.CreateTransactionDTO{
			AccountID: 123,
			Amount:    100.0,
			Direction: "credit",
		}
		account := &models.UserAccount{
			ID:      123,
			Balance: decimal.NewFromFloat(1000.0),
		}
		transaction := &models.Transaction{
			AccountID: 123,
			Amount:    100.0,
			Direction: "credit",
			Status:    "pending",
		}

		mockRepo.On("FindAccountById", 123).Return(account).Once()
		mockRepo.On("CreateTransaction", &transactionDTO).Return(transaction, nil)
		mockExternal.On("ForwardTransactionToThirdParty", transaction).Return(nil)
		mockRepo.On("UpdateTransactionStatus", transaction, "success").Return(nil)

		body, _ := json.Marshal(transactionDTO)
		req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockRepo.AssertExpectations(t)
		mockExternal.AssertExpectations(t)
	})

	t.Run("invalid payload", func(t *testing.T) {

		req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer([]byte("invalid payload")))
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("invalid amount", func(t *testing.T) {
		t.Skip()

		transactionDTO := dto.CreateTransactionDTO{
			AccountID: 123,
			Amount:    -100.0,
			Direction: "credit",
		}

		body, _ := json.Marshal(transactionDTO)
		req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("invalid account ID", func(t *testing.T) {
		t.Skip()
		transactionDTO := dto.CreateTransactionDTO{
			AccountID: 999,
			Amount:    100.0,
			Direction: "credit",
		}

		mockRepo.On("FindAccountById", 999).Return(nil).Once()

		body, _ := json.Marshal(transactionDTO)
		req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("forward transaction error", func(t *testing.T) {
		t.Skip()

		transactionDTO := dto.CreateTransactionDTO{
			AccountID: 123,
			Amount:    100.0,
			Direction: "credit",
		}
		account := &models.UserAccount{
			ID:      123,
			Balance: decimal.NewFromFloat(1000.0),
		}
		transaction := &models.Transaction{
			AccountID: 123,
			Amount:    100.0,
			Direction: "credit",
			Status:    "pending",
		}

		mockRepo.On("FindAccountById", 123).Return(account).Once()
		mockRepo.On("CreateTransaction", &transactionDTO).Return(transaction, nil).Once()
		mockExternal.On("ForwardTransactionToThirdParty", transaction).Return(errors.New("external error")).Once()
		mockRepo.On("UpdateTransactionStatus", transaction, "failed").Return(nil).Once()

		body, _ := json.Marshal(transactionDTO)
		req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockRepo.AssertExpectations(t)
		mockExternal.AssertExpectations(t)
	})
}
