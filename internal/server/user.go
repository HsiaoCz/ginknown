package server

import (
	"github.com/HsiaoCz/ginknown/types"
	"gorm.io/gorm"
)

type UserRepo interface {
	UserSingup(*types.User) error
	GetUserEmail(email string) bool
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

func (uc *UserUseCase) GetUserEmail(email string) bool {
	return true
}
