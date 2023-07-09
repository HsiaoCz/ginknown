package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/HsiaoCz/ginknown/etc"
	"github.com/HsiaoCz/ginknown/internal/service"
	"github.com/HsiaoCz/ginknown/storage"
	"github.com/gin-gonic/gin"
)

type Server struct {
	addr string
	port string
	r    *gin.Engine
	uc   *service.UserCase
}

func NewServer(r *gin.Engine, store *storage.Storage) *Server {
	return &Server{
		addr: etc.Conf.AC.Addr,
		port: etc.Conf.AC.Port,
		r:    r,
		uc:   service.NewUserCase(r, store),
	}
}

func (s *Server) Start() error {
	s.r.Use(gin.Recovery(), gin.Logger())
	s.uc.RegisterRouter()
	srv := http.Server{
		Handler:      s.r,
		Addr:         fmt.Sprintf("%s:%s", s.addr, s.port),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	return srv.ListenAndServe()
}
