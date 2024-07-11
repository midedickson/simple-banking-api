package mock_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/midedickson/simple-banking-app/dto"
	"github.com/midedickson/simple-banking-app/repository"
)

type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func CreateNewPOSTMockClient() *MockClient {
	return &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			var forwardTransactionDTO dto.ForwardTransactionDTO
			err := json.NewDecoder(req.Body).Decode(&forwardTransactionDTO)
			if err != nil {
				return &http.Response{StatusCode: http.StatusBadRequest}, err
			}
			repository.ExternalTransactions = append(repository.ExternalTransactions, &forwardTransactionDTO)
			return &http.Response{StatusCode: http.StatusOK, Body: req.Body}, nil
		},
	}
}

func CreateNewGETMockClient() *MockClient {
	return &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			urlPath := req.URL.Path
			parts := strings.Split(urlPath, "/")
			if len(parts) != 2 {
				return &http.Response{StatusCode: http.StatusBadRequest}, errors.New("invalid request path")
			}
			reference := parts[1]
			for _, transaction := range repository.ExternalTransactions {
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
