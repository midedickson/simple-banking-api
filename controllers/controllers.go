package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/midedickson/simple-banking-app/dto"
	"github.com/midedickson/simple-banking-app/external"
	"github.com/midedickson/simple-banking-app/repository"
	"github.com/midedickson/simple-banking-app/utils"
)

type controller struct {
	repo *repository.Repository
}

func NewController(repo *repository.Repository) *controller {
	return &controller{repo: repo}
}

func (c *controller) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var createTransactionDTO dto.CreateTransactionDTO
	err := json.NewDecoder(r.Body).Decode(&createTransactionDTO)
	if err != nil {
		utils.Dispatch400Error(w, "Invalid request payload", err)
		return
	}

	userAccount := c.repo.FindAccountById(createTransactionDTO.AccountID)
	if userAccount == nil {
		utils.Dispatch400Error(w, "Invalid account ID", nil)
		return
	}

	transaction, err := c.repo.NewTransaction(&createTransactionDTO)
	if err != nil {
		utils.Dispatch500Error(w, err)
		return
	}
	// todo: send transaction to the third-party system
	err = external.ForwardTransactionToThirdParty(transaction)
	if err != nil {
		c.repo.UpdateTransactionStatus(transaction, "failed")
		utils.Dispatch500Error(w, err)
		return
	}
	if transaction.Direction == "debit" {
		err := userAccount.Debit(transaction.Amount)
		if err != nil {
			c.repo.UpdateTransactionStatus(transaction, "failed")
			utils.Dispatch400Error(w, "Insufficient funds", err)
			return
		}
	} else {
		userAccount.Credit(transaction.Amount)
	}
	c.repo.UpdateTransactionStatus(transaction, "success")

	utils.Dispatch200(w, "Transaction created successfully", transaction)
}
func (c *controller) FetchTransactionDetails(w http.ResponseWriter, r *http.Request) {}
func (c *controller) FetchUserAccountDetails(w http.ResponseWriter, r *http.Request) {}
func (c *controller) Hello(w http.ResponseWriter, r *http.Request) {
	utils.Dispatch200(w, "hello, you have reached simple banking api", nil)
}
