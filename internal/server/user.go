package server

import (
	"github.com/HsiaoCz/ginknown/types"
	"gorm.io/gorm"
)

type UserRepo interface {
	UserSingup(*types.User) error
}

type UserUseCase struct {
	db *gorm.DB
}

func NewUserUseCase(db *gorm.DB) *UserUseCase {
	return &UserUseCase{
		db: db,
	}
}

func (uc *UserUseCase) UserSingup(user *types.User) (err error) {
	return nil
}
