package mocks

import "github.com/stretchr/testify/mock"

// Mock the repository
type MockIdempotencyStore struct {
	mock.Mock
}

func (m *MockIdempotencyStore) CreateNewIdempotencyKey() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockIdempotencyStore) CheckIdempotencyKeyStatus(key string) (string, error) {
	args := m.Called(key)

	return args.Get(0).(string), args.Error(1)
}

func (m *MockIdempotencyStore) UpdateIdempotencyKeyStatus(key string, status string) error {
	args := m.Called(key, status)
	return args.Error(0)
}
