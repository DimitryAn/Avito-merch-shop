package models

type User struct {
	Balance  int
	ID       int
	Username string `json:"username"`
	Password string `json:"password"`
}
