package idempotency

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/midedickson/simple-banking-app/constants"
)

func (s *KeyBasedIdempotencyStore) generateIdempotencyKey() (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// generate an idempotency key
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		newUUIDKey, err := uuid.NewUUID()
		if err != nil {
			return "", err
		}
		newStrKey := newUUIDKey.String()
		if _, ok := s.keyTable[newStrKey]; !ok {
			return newStrKey, nil
		}
	}
	return "", fmt.Errorf("failed to generate unique idempotency key after %d attempts", maxRetries)
}

func (s *KeyBasedIdempotencyStore) CreateNewIdempotencyKey() (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	key, err := s.generateIdempotencyKey()
	if err != nil {
		return "", err
	}
	s.keyTable[key] = constants.WAITING
	return key, nil
}

func (s *KeyBasedIdempotencyStore) ConfirmIdempotencyKeyAsProcessed(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.keyTable, key)
}

func (s *KeyBasedIdempotencyStore) CheckIdempotencyKeyStatus(key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	status, ok := s.keyTable[key]
	if !ok {
		return "", fmt.Errorf("requested idempotency key %v not found", key)
	}
	return status, nil
}

func (s *KeyBasedIdempotencyStore) UpdateIdempotencyKeyStatus(key string, status string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.keyTable[key]; ok {
		s.keyTable[key] = status
	} else {
		return fmt.Errorf("requested idempotency key for update %v not found", key)
	}
	return nil
}
