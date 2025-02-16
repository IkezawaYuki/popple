package entity

type User struct {
	Email    string `json:"string" binding:"required"`
	Password string `json:"password" binding:"required"`
}
