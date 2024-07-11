package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type controller struct {
	repo *repository
}

func newController(repo *repository) *controller {
	return &controller{repo: repo}
}

func (c *controller) createTransaction(w http.ResponseWriter, r *http.Request) {
	var createTransactionDTO createTransactionDTO
	err := json.NewDecoder(r.Body).Decode(&createTransactionDTO)
	if err != nil {
		Dispatch400Error(w, "Invalid request payload", err)
		return
	}

	userAccount := c.repo.findAccountById(createTransactionDTO.AccountID)
	if userAccount == nil {
		Dispatch400Error(w, "Invalid account ID", nil)
		return
	}

	transaction, err := c.repo.newTransaction(&createTransactionDTO)
	if err != nil {
		Dispatch500Error(w, err)
		return
	}
	// todo: send transaction to the third-party system
	err = forwardTransactionToThirdParty(transaction)
	if err != nil {
		c.repo.updateTransactionStatus(transaction, "failed")
		Dispatch500Error(w, err)
		return
	}
	if transaction.Direction == "debit" {
		err := userAccount.performAccountDebit(transaction.Amount)
		if err != nil {
			c.repo.updateTransactionStatus(transaction, "failed")
			Dispatch400Error(w, "Insufficient funds", err)
			return
		}
	} else {
		userAccount.performAccountCredit(transaction.Amount)
	}
	c.repo.updateTransactionStatus(transaction, "success")

	Dispatch200(w, "Transaction created successfully", transaction)
}
func (c *controller) fetchTransactionDetails(w http.ResponseWriter, r *http.Request) {}
func (c *controller) fetchUserAccountDetails(w http.ResponseWriter, r *http.Request) {}
func (c *controller) hello(w http.ResponseWriter, r *http.Request) {
	addDefaultHeaders(w)
	Dispatch200(w, "hello, you have reached simple banking api", nil)
}

// 500 - internal server error
func Dispatch500Error(w http.ResponseWriter, err error) {
	addDefaultHeaders(w)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(WriteError(fmt.Sprintf("%v", err), nil))
}

// 400 - bad request
func Dispatch400Error(w http.ResponseWriter, msg string, err any) {
	addDefaultHeaders(w)
	w.WriteHeader(http.StatusBadRequest)
	w.Write(WriteError(msg, err))
}

// 403 - forbidden request, incase of non-authorised request
func Dispatch403Error(w http.ResponseWriter, msg string, err any) {
	addDefaultHeaders(w)
	w.WriteHeader(http.StatusForbidden)
	w.Write(WriteError(msg, err))
}

// 404 - not found
func Dispatch404Error(w http.ResponseWriter, msg string, err any) {
	addDefaultHeaders(w)
	w.WriteHeader(http.StatusNotFound)
	w.Write(WriteError(msg, err))
}

// 200 - OK
func Dispatch200(w http.ResponseWriter, msg string, data any) {
	addDefaultHeaders(w)
	w.WriteHeader(http.StatusOK)
	w.Write(WriteInfo(msg, data))
}

func WriteInfo(message string, data any) []byte {
	response := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	r, err := json.Marshal(response)
	if err == nil {
		return r
	} else {
		log.Printf("err: %s", err)
	}
	return nil
}

func WriteError(message string, err interface{}) []byte {
	response := APIResponse{
		Success: false,
		Message: message,
		Data:    err,
	}
	data, err := json.Marshal(response)
	if err == nil {
		return data
	} else {
		log.Printf("err: %s", err)
	}
	return nil
}

func addDefaultHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
}

// get a path param from request
func GetPathParam(r *http.Request, name string) (string, error) {
	vars := mux.Vars(r)
	value, ok := vars[name]
	if !ok {
		return "", fmt.Errorf("invalid or missing %s in request param", name)
	}
	return value, nil

}
