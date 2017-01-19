package persistence

import (
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
)

// A persistence layer for Superside, using the filesystem as a backing store
type FileStore struct {
	basePath string
}

func NewFileStore(path string) *FileStore {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Infof("No persistence directory found, creating %s", path)
		os.MkdirAll(path, os.ModePerm)
	}
	return &FileStore{path}
}

func (r *FileStore) pathForKey(key string) string {
	return r.basePath + "/" + key + ".json"
}

func (r *FileStore) StoreBlob(key string, data []byte) error {
	err := ioutil.WriteFile(r.pathForKey(key), data, 0644)
	if err != nil {
		log.Warnf("Unable to write %s to store %s, %s", key, r.basePath, err)
	}
	return err
}


func (r *FileStore) GetBlob(key string) ([]byte, error) {
	if _, err := os.Stat(r.pathForKey(key)); os.IsNotExist(err) {
		log.Warn("No persistence file found, skipping")
		return []byte{}, nil
	}
	return ioutil.ReadFile(r.pathForKey(key))
}
