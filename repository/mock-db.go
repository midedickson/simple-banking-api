package repository

import (
	"github.com/midedickson/simple-banking-app/dto"
	"github.com/midedickson/simple-banking-app/models"
)

// fake db of users
var Users []*models.UserAccount = []*models.UserAccount{
	{AccountID: 1, Balance: 100},
	{AccountID: 2, Balance: 300},
	{AccountID: 3, Balance: 400},
}

// fake db of external trasnactions
var ExternalTransactions []*dto.ForwardTransactionDTO = []*dto.ForwardTransactionDTO{}
