package main

// fake db of users
var users []*UserAccount = []*UserAccount{
	{AccountID: 1, Balance: 100},
	{AccountID: 2, Balance: 300},
	{AccountID: 3, Balance: 400},
}

// fake db of external trasnactions
var externalTransactions []*forwardTransactionDTO = []*forwardTransactionDTO{}
