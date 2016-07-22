package persistence

type Store interface {
	StoreBlob(key string, data []byte) error
	GetBlob(key string) ([]byte, error)
}
