package idempotency

type IdempotencyStore interface {
	CreateNewIdempotencyKey() (string, error)
	CheckIdempotencyKeyStatus(key string) (string, error)
	UpdateIdempotencyKeyStatus(key string, status string) error
}

type KeyBasedIdempotencyStore struct {
	keyTable map[string]string
}

func NewIdempotencyStore() *KeyBasedIdempotencyStore {
	return &KeyBasedIdempotencyStore{keyTable: make(map[string]string)}
}
