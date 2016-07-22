package persistence

import (
	"gopkg.in/redis.v4"
)

// A persistence layer for Superside, using Redis as a backing store
type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(addr string, password string, db int) (*RedisStore, error){
	client := redis.NewClient(&redis.Options{
		Addr: addr,
		Password: password,
		DB: db,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &RedisStore{ client: client }, nil
}

func (r *RedisStore) StoreBlob(key string, data []byte) error {
	_, err := r.client.Set(key, data, 0).Result()
	return err
}

func (r *RedisStore) GetBlob(key string) ([]byte, error) {
	data, err := r.client.Get(key).Result()
	if err != nil {
		return nil, err
	}

	return []byte(data), nil
}
