package repository

import (
	"github.com/midedickson/simple-banking-app/dto"
	"github.com/midedickson/simple-banking-app/models"
	"github.com/shopspring/decimal"
)

// fake db of users
var Users []*models.UserAccount = []*models.UserAccount{
	{ID: 1, Balance: decimal.NewFromFloat(400.0)},
	{ID: 2, Balance: decimal.NewFromFloat(400.0)},
	{ID: 3, Balance: decimal.NewFromFloat(400.0)},
}

// fake db of external trasnactions
var ExternalTransactions []*dto.ForwardTransactionDTO = []*dto.ForwardTransactionDTO{}
