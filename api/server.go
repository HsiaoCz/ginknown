package api

import (
	"net/http"
	"time"

	"github.com/HsiaoCz/ginknown/internal/service"
	"github.com/HsiaoCz/ginknown/storage"
	"github.com/gin-gonic/gin"
)

type Server struct {
	r  *gin.Engine
	uc *service.UserCase
}

func NewServer(r *gin.Engine, store *storage.Storage) *Server {
	return &Server{
		r:  r,
		uc: service.NewUserCase(r, store),
	}
}

func (s *Server) Start() error {
	s.r.Use(gin.Recovery(), gin.Logger())
	srv := http.Server{
		Handler:      s.r,
		Addr:         ":9091",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	return srv.ListenAndServe()
}
