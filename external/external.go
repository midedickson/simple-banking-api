package external

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/midedickson/simple-banking-app/dto"
	mock_client "github.com/midedickson/simple-banking-app/mock"
	"github.com/midedickson/simple-banking-app/models"
)

var ErrThirdPartyFailure = errors.New("third-party failure")

func FetchTransactionDetailsFromThirdParty(reference string) (*dto.ForwardTransactionDTO, error) {
	var transaction *dto.ForwardTransactionDTO
	client := mock_client.CreateNewGETMockClient()
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

func ForwardTransactionToThirdParty(transaction *models.Transaction) error {
	client := mock_client.CreateNewPOSTMockClient()
	forwardTransactionDto := &dto.ForwardTransactionDTO{
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
