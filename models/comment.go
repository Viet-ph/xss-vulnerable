package models

type Comment struct {
	Body string `json:"body"`
	Id   int    `json:"id"`
}
