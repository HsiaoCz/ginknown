package storage

type RedisStorage struct{}

func NewRedisStorage() *RedisStorage {
	return &RedisStorage{}
}

func (rs *RedisStorage) InitStorage() error {
	return nil
}
