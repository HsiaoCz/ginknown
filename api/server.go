package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	r *gin.Engine
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Start() error {
	s.r = gin.Default()

	srv := http.Server{
		Handler:      s.r,
		Addr:         ":9091",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	return srv.ListenAndServe()
}
