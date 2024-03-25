package models

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Id       int    `json:"id"`
}
