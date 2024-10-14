package idempotency

import "sync"

type IdempotencyStore interface {
	CreateNewIdempotencyKey() (string, error)
	CheckIdempotencyKeyStatus(key string) (string, error)
	UpdateIdempotencyKeyStatus(key string, status string) error
}

type KeyBasedIdempotencyStore struct {
	keyTable map[string]string
	mu       sync.Mutex
}

func NewIdempotencyStore() *KeyBasedIdempotencyStore {
	return &KeyBasedIdempotencyStore{keyTable: make(map[string]string), mu: sync.Mutex{}}
}
