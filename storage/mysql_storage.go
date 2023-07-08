package storage

type MysqlStorage struct {
}

func NewMysqlStorage() *MysqlStorage {
	return &MysqlStorage{}
}

func (ms *MysqlStorage) InitStorage() error {
	return nil
}
