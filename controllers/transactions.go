package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/midedickson/simple-banking-app/constants"
	"github.com/midedickson/simple-banking-app/dto"
	"github.com/midedickson/simple-banking-app/external"
	"github.com/midedickson/simple-banking-app/idempotency"
	"github.com/midedickson/simple-banking-app/repository"
	"github.com/midedickson/simple-banking-app/utils"
	"github.com/shopspring/decimal"
)

type Controller struct {
	repo             repository.Repository
	external         external.External
	idempotencyStore idempotency.IdempotencyStore
}

func NewController(repo repository.Repository, external external.External, idempotencyStore idempotency.IdempotencyStore) *Controller {
	return &Controller{repo: repo, external: external, idempotencyStore: idempotencyStore}
}

// func (c *Controller) CheckIdempotencyKeyStatus(key string) (string, error) {

// 	return status, nil
// }

func (c *Controller) CreateCreditTransaction(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("X-Idempotency-Key")
	if key == "" {
		utils.Dispatch400Error(w, "Idempotency Key is required", nil)
		return
	}
	status, err := c.idempotencyStore.CheckIdempotencyKeyStatus(key)
	if err != nil {
		utils.Dispatch400Error(w, "Idempotency Key is required", nil)
		return
	}
	switch status {
	case constants.SUCCESS:
		utils.Dispatch409Error(w, "Idempotency Key has already been processed", status)
		return
	case constants.WAITING:
		c.idempotencyStore.UpdateIdempotencyKeyStatus(key, constants.PROCESSING)
	case constants.PROCESSING:
		utils.Dispatch409Error(w, "A similar is already being processed, please wait to get the a feedback and try again later if it doesn't work.", status)
		return
	case constants.FAILED:
		utils.Dispatch500Error(w, errors.New("a similar transaction has failed, please try again"))
		return
	}
	var createTransactionDTO dto.CreateTransactionDTO
	err = json.NewDecoder(r.Body).Decode(&createTransactionDTO)
	if err != nil {
		c.idempotencyStore.UpdateIdempotencyKeyStatus(key, constants.FAILED)
		utils.Dispatch400Error(w, "Invalid request payload", err)
		return
	}
	zero := decimal.NewFromInt(0)
	amountToAdd := decimal.NewFromFloat(createTransactionDTO.Amount)

	if amountToAdd.LessThanOrEqual(zero) {
		c.idempotencyStore.UpdateIdempotencyKeyStatus(key, constants.FAILED)

		utils.Dispatch400Error(w, "Invalid Amount", nil)
		return
	}
	userAccount := c.repo.FindAccountById(createTransactionDTO.AccountID)
	if userAccount == nil {
		c.idempotencyStore.UpdateIdempotencyKeyStatus(key, constants.FAILED)

		utils.Dispatch400Error(w, "Invalid account ID", nil)
		return
	}
	createDBTransactionDTO := &dto.CreateDBTransactionDTO{
		AccountID: createTransactionDTO.AccountID,
		Amount:    createTransactionDTO.Amount,
		Direction: constants.DirectionCredit,
	}

	transaction, err := c.repo.CreateTransaction(createDBTransactionDTO)
	if err != nil {
		c.idempotencyStore.UpdateIdempotencyKeyStatus(key, constants.FAILED)
		utils.Dispatch500Error(w, err)
		return
	}
	// send transaction to the third-party system
	err = c.external.ForwardTransactionToThirdParty(transaction)
	if err != nil {
		c.idempotencyStore.UpdateIdempotencyKeyStatus(key, constants.FAILED)
		c.repo.UpdateTransactionStatus(transaction, constants.FAILED)
		utils.Dispatch500Error(w, err)
		return
	}

	err = userAccount.Credit(transaction.Amount)
	if err != nil {
		c.idempotencyStore.UpdateIdempotencyKeyStatus(key, constants.FAILED)
		c.repo.UpdateTransactionStatus(transaction, constants.FAILED)
		utils.Dispatch500Error(w, err)
		return
	}
	c.idempotencyStore.UpdateIdempotencyKeyStatus(key, constants.SUCCESS)
	c.repo.UpdateTransactionStatus(transaction, constants.SUCCESS)

	utils.Dispatch200(w, "Transaction created successfully", transaction)
}

func (c *Controller) CreateDebitTransaction(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("X-Idempotency-Key")

	status, err := c.idempotencyStore.CheckIdempotencyKeyStatus(key)
	if err != nil {
		utils.Dispatch422Error(w, "Invalid or missing Idempotency Key", nil)
		return
	}
	switch status {
	case constants.SUCCESS:
		utils.Dispatch409Error(w, "Idempotency Key has already been processed", status)
		return
	case constants.WAITING:
		c.idempotencyStore.UpdateIdempotencyKeyStatus(key, constants.PROCESSING)
	case constants.PROCESSING:
		utils.Dispatch409Error(w, "A similar transaction is already being processed, please wait to get the a feedback and try again later if it doesn't work.", status)
		return
	case constants.FAILED:
		utils.Dispatch409Error(w, "A similar transaction has failed, please try again.", status)
		return
	}
	var createTransactionDTO dto.CreateTransactionDTO
	err = json.NewDecoder(r.Body).Decode(&createTransactionDTO)
	if err != nil {
		c.idempotencyStore.UpdateIdempotencyKeyStatus(key, constants.FAILED)
		utils.Dispatch400Error(w, "Invalid request payload", err)
		return
	}
	zero := decimal.NewFromInt(0)
	amountToAdd := decimal.NewFromFloat(createTransactionDTO.Amount)

	if amountToAdd.LessThanOrEqual(zero) {
		c.idempotencyStore.UpdateIdempotencyKeyStatus(key, constants.FAILED)
		utils.Dispatch400Error(w, "Invalid Amount", nil)
		return
	}
	userAccount := c.repo.FindAccountById(createTransactionDTO.AccountID)
	if userAccount == nil {
		c.idempotencyStore.UpdateIdempotencyKeyStatus(key, constants.FAILED)
		utils.Dispatch400Error(w, "Invalid account ID", nil)
		return
	}

	createDBTransactionDTO := &dto.CreateDBTransactionDTO{
		AccountID: createTransactionDTO.AccountID,
		Amount:    createTransactionDTO.Amount,
		Direction: constants.DirectionDebit,
	}

	transaction, err := c.repo.CreateTransaction(createDBTransactionDTO)
	if err != nil {
		utils.Dispatch500Error(w, err)
		return
	}
	// send transaction to the third-party system
	err = c.external.ForwardTransactionToThirdParty(transaction)
	if err != nil {
		c.idempotencyStore.UpdateIdempotencyKeyStatus(key, constants.FAILED)
		c.repo.UpdateTransactionStatus(transaction, constants.FAILED)
		utils.Dispatch500Error(w, err)
		return
	}
	err = userAccount.Debit(transaction.Amount)
	if err != nil {
		c.idempotencyStore.UpdateIdempotencyKeyStatus(key, constants.FAILED)
		c.repo.UpdateTransactionStatus(transaction, constants.FAILED)
		utils.Dispatch400Error(w, "Insufficient funds", err)
		return
	}
	c.idempotencyStore.UpdateIdempotencyKeyStatus(key, constants.SUCCESS)
	c.repo.UpdateTransactionStatus(transaction, constants.SUCCESS)

	utils.Dispatch200(w, "Transaction created successfully", transaction)
}
func (c *Controller) FetchTransactionDetails(w http.ResponseWriter, r *http.Request) {}
func (c *Controller) FetchUserAccountDetails(w http.ResponseWriter, r *http.Request) {}
func (c *Controller) Hello(w http.ResponseWriter, r *http.Request) {
	utils.Dispatch200(w, "hello, you have reached simple banking api", nil)
}
