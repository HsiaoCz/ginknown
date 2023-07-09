package main

import (
	"log"

	"github.com/HsiaoCz/ginknown/api"
	"github.com/HsiaoCz/ginknown/etc"
	"github.com/HsiaoCz/ginknown/storage"
	"github.com/gin-gonic/gin"
)

func main() {

	if err := etc.InitConf(); err != nil {
		log.Fatal(err)
	}

	if err := storage.NewStorage().Is.StartConn(storage.NewMysqlStorage(), storage.NewRedisStorage()); err != nil {
		log.Fatal(err)
	}

	if err := api.NewServer(gin.New(), storage.NewStorage()).Start(); err != nil {
		log.Fatal(err)
	}

}
