package persistence

import (
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
)

// A persistence layer for Superside, using Redis as a backing store
type FileStore struct {
	basePath string
}

func NewFileStore(path string) *FileStore {
	return &FileStore{path}
}

func (r *FileStore) pathForKey(key string) string {
	return r.basePath + "/" + key + ".json"
}

func (r *FileStore) StoreBlob(key string, data []byte) error {
	return ioutil.WriteFile(r.pathForKey(key), data, 0644)
}


func (r *FileStore) GetBlob(key string) ([]byte, error) {
	if _, err := os.Stat(r.pathForKey(key)); os.IsNotExist(err) {
		log.Warn("No persistence file found, skipping")
		return []byte{}, nil
	}
	return ioutil.ReadFile(r.pathForKey(key))
}
