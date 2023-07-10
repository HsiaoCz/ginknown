package service

import (
	"net/http"

	"github.com/HsiaoCz/ginknown/storage"
	"github.com/HsiaoCz/ginknown/types"
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

func (u *UserCase) UserSingup(c *gin.Context) {
	s := new(types.UserRegister)
	err := c.ShouldBind(s)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"Code":    4001,
			"Message": "参数错误,请检查",
		})
		return
	}

	isHasEmail := u.store.Ms.UR.GetUserEmail(s.Email)
	if isHasEmail {
		c.JSON(http.StatusOK, gin.H{
			"Code":    2001,
			"Message": "该邮箱已经注册,请切换邮箱",
		})
		return
	}

}
