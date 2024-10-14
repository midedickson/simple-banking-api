package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/midedickson/simple-banking-app/constants"
	"github.com/midedickson/simple-banking-app/controllers"
	"github.com/midedickson/simple-banking-app/dto"
	"github.com/midedickson/simple-banking-app/models"
	"github.com/midedickson/simple-banking-app/tests/mocks"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCreateDebitTransaction(t *testing.T) {
	mockRepo := new(mocks.MockRepo)
	mockExternal := new(mocks.MockExternal)
	mockIdempotencyStore := new(mocks.MockIdempotencyStore)
	ctrl := controllers.NewController(mockRepo, mockExternal, mockIdempotencyStore)
	handler := http.HandlerFunc(ctrl.CreateDebitTransaction)

	t.Run("successful debit transaction", func(t *testing.T) {
		transactionDTO := dto.CreateTransactionDTO{
			AccountID: 123,
			Amount:    100.0,
		}
		account := &models.UserAccount{
			ID:      123,
			Balance: decimal.NewFromFloat(1000.0),
		}
		createDBTransactionDTO := dto.CreateDBTransactionDTO{
			AccountID: transactionDTO.AccountID,
			Amount:    transactionDTO.Amount,
			Direction: "debit",
		}
		transaction := &models.Transaction{
			AccountID: 123,
			Amount:    100.0,
			Direction: "debit",
			Status:    "pending",
		}

		mockIdempotencyStore.On("CheckIdempotencyKeyStatus", "12345").Return(constants.WAITING, nil)
		mockIdempotencyStore.On("UpdateIdempotencyKeyStatus", "12345", constants.PROCESSING).Return(nil)
		mockRepo.On("FindAccountById", 123).Return(account)
		mockRepo.On("CreateTransaction", &createDBTransactionDTO).Return(transaction, nil)
		mockExternal.On("ForwardTransactionToThirdParty", transaction).Return(nil)
		mockRepo.On("UpdateTransactionStatus", transaction, "success").Return(nil)
		mockIdempotencyStore.On("UpdateIdempotencyKeyStatus", "12345", constants.SUCCESS).Return(nil)

		body, _ := json.Marshal(transactionDTO)
		req, _ := http.NewRequest("POST", "/transactions/debit", bytes.NewBuffer(body))
		req.Header.Set("X-Idempotency-Key", "12345")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockRepo.AssertExpectations(t)
		mockExternal.AssertExpectations(t)
		mockIdempotencyStore.AssertExpectations(t)
	})

	t.Run("forward debit transaction error", func(t *testing.T) {
		transactionDTO := dto.CreateTransactionDTO{
			AccountID: 124,
			Amount:    100.0,
		}
		account := &models.UserAccount{
			ID:      124,
			Balance: decimal.NewFromFloat(1000.0),
		}
		createDBTransactionDTO := dto.CreateDBTransactionDTO{
			AccountID: transactionDTO.AccountID,
			Amount:    transactionDTO.Amount,
			Direction: "debit",
		}
		transaction := &models.Transaction{
			AccountID: createDBTransactionDTO.AccountID,
			Amount:    createDBTransactionDTO.Amount,
			Direction: "debit",
			Status:    "pending",
		}
		mockIdempotencyStore.On("CheckIdempotencyKeyStatus", "12345").Return(constants.WAITING, nil)
		mockIdempotencyStore.On("UpdateIdempotencyKeyStatus", "12345", constants.PROCESSING).Return(nil)

		mockRepo.On("FindAccountById", 124).Return(account)
		mockRepo.On("CreateTransaction", &createDBTransactionDTO).Return(transaction, nil)
		mockExternal.On("ForwardTransactionToThirdParty", transaction).Return(constants.ErrThirdPartyFailure)
		mockRepo.On("UpdateTransactionStatus", transaction, "failed").Return(nil).Once()
		mockIdempotencyStore.On("UpdateIdempotencyKeyStatus", "12345", constants.FAILED).Return(nil)

		body, _ := json.Marshal(transactionDTO)
		req, _ := http.NewRequest("POST", "/transactions/debit", bytes.NewBuffer(body))
		req.Header.Set("X-Idempotency-Key", "12345")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockRepo.AssertExpectations(t)
		mockExternal.AssertExpectations(t)
		mockIdempotencyStore.AssertExpectations(t)
	})

	t.Run("create pending transaction failed", func(t *testing.T) {

		transactionDTO := dto.CreateTransactionDTO{
			AccountID: 125,
			Amount:    100.0,
		}
		account := &models.UserAccount{
			ID:      125,
			Balance: decimal.NewFromFloat(1000.0),
		}
		createDBTransactionDTO := &dto.CreateDBTransactionDTO{
			AccountID: transactionDTO.AccountID,
			Amount:    transactionDTO.Amount,
			Direction: "debit",
		}
		transaction := &models.Transaction{
			AccountID: createDBTransactionDTO.AccountID,
			Amount:    createDBTransactionDTO.Amount,
			Direction: "debit",
			Status:    "pending",
		}
		mockIdempotencyStore.On("CheckIdempotencyKeyStatus", "12345").Return(constants.WAITING, nil)
		mockIdempotencyStore.On("UpdateIdempotencyKeyStatus", "12345", constants.PROCESSING).Return(nil)

		mockRepo.On("FindAccountById", 125).Return(account)
		mockRepo.On("CreateTransaction", createDBTransactionDTO).Return(transaction, errors.New("error creating transaction in DB"))
		mockIdempotencyStore.On("UpdateIdempotencyKeyStatus", "12345", constants.FAILED).Return(nil)

		body, _ := json.Marshal(transactionDTO)
		req, _ := http.NewRequest("POST", "/transactions/debit", bytes.NewBuffer(body))
		req.Header.Set("X-Idempotency-Key", "12345")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockRepo.AssertExpectations(t)
		mockIdempotencyStore.AssertExpectations(t)
	})
}

func TestCreateCreditTransaction(t *testing.T) {
	mockRepo := new(mocks.MockRepo)
	mockExternal := new(mocks.MockExternal)
	mockIdempotencyStore := new(mocks.MockIdempotencyStore)
	ctrl := controllers.NewController(mockRepo, mockExternal, mockIdempotencyStore)

	handler := http.HandlerFunc(ctrl.CreateCreditTransaction)

	t.Run("successful credit transaction", func(t *testing.T) {
		transactionDTO := dto.CreateTransactionDTO{
			AccountID: 123,
			Amount:    100.0,
		}
		account := &models.UserAccount{
			ID:      123,
			Balance: decimal.NewFromFloat(1000.0),
		}
		createDBTransactionDTO := dto.CreateDBTransactionDTO{
			AccountID: transactionDTO.AccountID,
			Amount:    transactionDTO.Amount,
			Direction: "credit",
		}
		transaction := &models.Transaction{
			AccountID: createDBTransactionDTO.AccountID,
			Amount:    createDBTransactionDTO.Amount,
			Direction: "credit",
			Status:    "pending",
		}
		mockIdempotencyStore.On("CheckIdempotencyKeyStatus", "12345").Return(constants.WAITING, nil)
		mockIdempotencyStore.On("UpdateIdempotencyKeyStatus", "12345", constants.PROCESSING).Return(nil)

		mockRepo.On("FindAccountById", 123).Return(account).Once()
		mockRepo.On("CreateTransaction", &createDBTransactionDTO).Return(transaction, nil)
		mockExternal.On("ForwardTransactionToThirdParty", transaction).Return(nil)
		mockRepo.On("UpdateTransactionStatus", transaction, "success").Return(nil)
		mockIdempotencyStore.On("UpdateIdempotencyKeyStatus", "12345", constants.SUCCESS).Return(nil)
		body, _ := json.Marshal(transactionDTO)
		req, _ := http.NewRequest("POST", "/transactions/credit", bytes.NewBuffer(body))
		req.Header.Set("X-Idempotency-Key", "12345")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockRepo.AssertExpectations(t)
		mockExternal.AssertExpectations(t)
		mockIdempotencyStore.AssertExpectations(t)
	})

	t.Run("invalid payload", func(t *testing.T) {
		mockIdempotencyStore.On("CheckIdempotencyKeyStatus", "12345").Return(constants.WAITING, nil)
		mockIdempotencyStore.On("UpdateIdempotencyKeyStatus", "12345", constants.PROCESSING).Return(nil)
		mockIdempotencyStore.On("UpdateIdempotencyKeyStatus", "12345", constants.FAILED).Return(nil)

		req, _ := http.NewRequest("POST", "/transactions/credit", bytes.NewBuffer([]byte("invalid payload")))
		req.Header.Set("X-Idempotency-Key", "12345")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockIdempotencyStore.AssertExpectations(t)
	})

	t.Run("invalid amount", func(t *testing.T) {

		transactionDTO := dto.CreateTransactionDTO{
			AccountID: 123,
			Amount:    -100.0,
		}
		mockIdempotencyStore.On("CheckIdempotencyKeyStatus", "12345").Return(constants.WAITING, nil)
		mockIdempotencyStore.On("UpdateIdempotencyKeyStatus", "12345", constants.PROCESSING).Return(nil)
		mockIdempotencyStore.On("UpdateIdempotencyKeyStatus", "12345", constants.FAILED).Return(nil)
		body, _ := json.Marshal(transactionDTO)
		req, _ := http.NewRequest("POST", "/transactions/credit", bytes.NewBuffer(body))
		req.Header.Set("X-Idempotency-Key", "12345")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockIdempotencyStore.AssertExpectations(t)

	})

	t.Run("invalid account ID", func(t *testing.T) {
		transactionDTO := dto.CreateTransactionDTO{
			AccountID: 999,
			Amount:    100.0,
		}

		mockRepo.On("FindAccountById", 999).Return(nil).Once()
		mockIdempotencyStore.On("CheckIdempotencyKeyStatus", "12345").Return(constants.WAITING, nil)
		mockIdempotencyStore.On("UpdateIdempotencyKeyStatus", "12345", constants.PROCESSING).Return(nil)
		mockIdempotencyStore.On("UpdateIdempotencyKeyStatus", "12345", constants.FAILED).Return(nil)
		body, _ := json.Marshal(transactionDTO)
		req, _ := http.NewRequest("POST", "/transactions/credit", bytes.NewBuffer(body))
		req.Header.Set("X-Idempotency-Key", "12345")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockRepo.AssertExpectations(t)
		mockIdempotencyStore.AssertExpectations(t)
	})

	t.Run("forward transaction error", func(t *testing.T) {

		transactionDTO := dto.CreateTransactionDTO{
			AccountID: 124,
			Amount:    100.0,
		}
		account := &models.UserAccount{
			ID:      124,
			Balance: decimal.NewFromFloat(1000.0),
		}
		createDBTransactionDTO := dto.CreateDBTransactionDTO{
			AccountID: transactionDTO.AccountID,
			Amount:    transactionDTO.Amount,
			Direction: "credit",
		}
		transaction := &models.Transaction{
			AccountID: createDBTransactionDTO.AccountID,
			Amount:    createDBTransactionDTO.Amount,
			Direction: "credit",
			Status:    "pending",
		}
		mockIdempotencyStore.On("CheckIdempotencyKeyStatus", "12345").Return(constants.WAITING, nil)
		mockIdempotencyStore.On("UpdateIdempotencyKeyStatus", "12345", constants.PROCESSING).Return(nil)
		mockRepo.On("FindAccountById", 124).Return(account)
		mockRepo.On("CreateTransaction", &createDBTransactionDTO).Return(transaction, nil)
		mockExternal.On("ForwardTransactionToThirdParty", transaction).Return(constants.ErrThirdPartyFailure)
		mockRepo.On("UpdateTransactionStatus", transaction, "failed").Return(nil)
		mockIdempotencyStore.On("UpdateIdempotencyKeyStatus", "12345", constants.FAILED).Return(nil)

		body, _ := json.Marshal(transactionDTO)
		req, _ := http.NewRequest("POST", "/transactions/credit", bytes.NewBuffer(body))
		req.Header.Set("X-Idempotency-Key", "12345")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockExternal.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
		mockIdempotencyStore.AssertExpectations(t)
	})
}

func TestProcessingIdempotency(t *testing.T) {
	mockRepo := new(mocks.MockRepo)
	mockExternal := new(mocks.MockExternal)
	mockIdempotencyStore := new(mocks.MockIdempotencyStore)
	ctrl := controllers.NewController(mockRepo, mockExternal, mockIdempotencyStore)

	handler := http.HandlerFunc(ctrl.CreateCreditTransaction)
	t.Run("idempotency processing", func(t *testing.T) {

		transactionDTO := dto.CreateTransactionDTO{
			AccountID: 124,
			Amount:    100.0,
		}

		mockIdempotencyStore.On("CheckIdempotencyKeyStatus", "12345").Return(constants.PROCESSING, nil)

		body, _ := json.Marshal(transactionDTO)
		req, _ := http.NewRequest("POST", "/transactions/credit", bytes.NewBuffer(body))
		req.Header.Set("X-Idempotency-Key", "12345")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusConflict, rr.Code)
		mockExternal.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
		mockIdempotencyStore.AssertExpectations(t)
	})
	t.Run("idempotency success", func(t *testing.T) {

		transactionDTO := dto.CreateTransactionDTO{
			AccountID: 124,
			Amount:    100.0,
		}

		mockIdempotencyStore.On("CheckIdempotencyKeyStatus", "12345").Return(constants.SUCCESS, nil)

		body, _ := json.Marshal(transactionDTO)
		req, _ := http.NewRequest("POST", "/transactions/credit", bytes.NewBuffer(body))
		req.Header.Set("X-Idempotency-Key", "12345")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusConflict, rr.Code)
		mockExternal.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
		mockIdempotencyStore.AssertExpectations(t)
	})

	t.Run("idempotency failed", func(t *testing.T) {

		transactionDTO := dto.CreateTransactionDTO{
			AccountID: 124,
			Amount:    100.0,
		}

		mockIdempotencyStore.On("CheckIdempotencyKeyStatus", "12345").Return(constants.FAILED, nil)

		body, _ := json.Marshal(transactionDTO)
		req, _ := http.NewRequest("POST", "/transactions/credit", bytes.NewBuffer(body))
		req.Header.Set("X-Idempotency-Key", "12345")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusConflict, rr.Code)
		mockExternal.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
		mockIdempotencyStore.AssertExpectations(t)
	})
}
