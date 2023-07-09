package storage

import (
	"fmt"

	"github.com/HsiaoCz/ginknown/etc"
	"github.com/HsiaoCz/ginknown/internal/server"
	"github.com/HsiaoCz/ginknown/types"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

// mysql conf
type mysqlConf struct {
	mysql_user     string
	mysql_password string
	mysql_host     string
	mysql_port     string
	db_name        string
}

type MysqlStorage struct {
	mc *mysqlConf
	ur server.UserRepo
}

func NewMysqlStorage() *MysqlStorage {
	mcy := etc.Conf.MC
	return &MysqlStorage{
		mc: &mysqlConf{
			mysql_user:     mcy.Mysql_User,
			mysql_password: mcy.Password,
			mysql_host:     mcy.Mysql_Host,
			mysql_port:     mcy.Mysql_Port,
			db_name:        mcy.DB_Name,
		},
		ur: server.NewUserUseCase(db),
	}
}

func (ms *MysqlStorage) InitStorage() error {
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", ms.mc.mysql_user, ms.mc.mysql_password, ms.mc.mysql_host, ms.mc.mysql_port, ms.mc.db_name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	db.AutoMigrate(&types.User{})
	return err
}
