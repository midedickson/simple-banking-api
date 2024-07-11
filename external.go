package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

var ErrThirdPartyFailure = errors.New("third-party failure")

type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func createNewPOSTMockClient() *MockClient {
	return &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			var forwardTransactionDTO forwardTransactionDTO
			err := json.NewDecoder(req.Body).Decode(&forwardTransactionDTO)
			if err != nil {
				return &http.Response{StatusCode: http.StatusBadRequest}, err
			}
			externalTransactions = append(externalTransactions, &forwardTransactionDTO)
			return &http.Response{StatusCode: http.StatusOK, Body: req.Body}, nil
		},
	}
}

func createNewGETMockClient() *MockClient {
	return &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			urlPath := req.URL.Path
			parts := strings.Split(urlPath, "/")
			if len(parts) != 2 {
				return &http.Response{StatusCode: http.StatusBadRequest}, errors.New("invalid request path")
			}
			reference := parts[1]
			for _, transaction := range externalTransactions {
				if transaction.Reference == reference {
					data, err := json.Marshal(transaction)
					if err != nil {
						log.Printf("failed to marshal transaction data from third party: %s", err)
						return &http.Response{StatusCode: http.StatusBadRequest}, err
					}
					return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(data))}, nil
				}
			}
			return &http.Response{StatusCode: http.StatusNotFound}, errors.New("transaction not found")
		},
	}
}

func fetchTransactionDetailsFromThirdParty(reference string) (*forwardTransactionDTO, error) {
	var transaction *forwardTransactionDTO
	client := createNewGETMockClient()
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://third-party-system.com/transactions/%s", reference), nil)
	if err != nil {
		log.Printf("failed to create request for third party: %s", err)
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("failed to send request to third party: %s", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to forward transaction to third party: %s", resp.Status)
		return nil, ErrThirdPartyFailure
	}
	log.Printf("Transaction forwarded successfully to third party: %+v", transaction)
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to parse transaction from third party: %s", err)
		return nil, err
	}
	err = json.Unmarshal(bodyBytes, &transaction)
	if err != nil {
		log.Printf("failed to parse transaction from third party: %s", err)
		return nil, err
	}

	return transaction, nil
}

func forwardTransactionToThirdParty(transaction *Transaction) error {
	client := createNewPOSTMockClient()
	forwardTransactionDto := &forwardTransactionDTO{
		Reference: transaction.Reference,
		AccountID: transaction.AccountID,
		Amount:    transaction.Amount,
	}
	data, err := json.Marshal(forwardTransactionDto)
	if err != nil {
		log.Printf("failed to marshal transaction data foer third party: %s", err)
		return err
	}
	req, err := http.NewRequest(http.MethodPost, "http://third-party-system.com/transactions", bytes.NewReader(data))
	if err != nil {
		log.Printf("failed to create request for third party: %s", err)
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("failed to send request to third party: %s", err)
		return err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to forward transaction to third party: %s", resp.Status)
		return ErrThirdPartyFailure
	}
	log.Printf("Transaction forwarded successfully to third party: %+v", transaction)
	return nil
}
