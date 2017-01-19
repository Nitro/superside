package persistence

type Store interface {
	StoreBlob(key string, data []byte) error
	GetBlob(key string) ([]byte, error)
}

type NoopStore struct{}

func (n *NoopStore) StoreBlob(key string, data []byte) error {
	return nil
}

func (n *NoopStore) GetBlob(key string) ([]byte, error) {
	return []byte{}, nil
}
