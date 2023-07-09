package service

import (
	"github.com/HsiaoCz/ginknown/storage"
	"github.com/gin-gonic/gin"
)

type UserCase struct {
	r     *gin.Engine
	store *storage.Storage
}

func NewUserCase(r *gin.Engine, store *storage.Storage) *UserCase {
	return &UserCase{
		r:     r,
		store: store,
	}
}

func (u *UserCase) RegisterRouter() {
	u.r.POST("/user/singup", u.UserSingup)
}

func (u *UserCase) UserSingup(c *gin.Context) {}
