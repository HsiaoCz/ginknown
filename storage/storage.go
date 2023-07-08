package storage

type Storage struct {
	Ms *MysqlStorage
	Rs *RedisStorage
	Is *IStore
}

func NewStorage() *Storage {
	return &Storage{
		Ms: NewMysqlStorage(),
		Rs: NewRedisStorage(),
		Is: NewIStore(),
	}
}
