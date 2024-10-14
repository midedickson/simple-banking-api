package mocks

import (
	"github.com/midedickson/simple-banking-app/models"
	"github.com/stretchr/testify/mock"
)

// Mock the external service
type MockExternal struct {
	mock.Mock
}

func (m *MockExternal) ForwardTransactionToThirdParty(transaction *models.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}
