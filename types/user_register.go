package types

type UserRegister struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"eqfiled=Password"`
	Email      string `json:"email" binding:"required"`
}
