package model

import (
	"devopscenter/helper"
)

type Login struct {
	Id       int    `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

func (l Login) SearchPassword(user string) (password string) {
	mysqlEngine := helper.SqlContext
	mysqlEngine.QueryRow("SELECT password FROM user WHERE username = ?", user).Scan(&password)
	return password
}

func (l Login) SearchUser(username string) (exists int) {
	mysqlEngine := helper.SqlContext
	mysqlEngine.QueryRow("SELECT EXISTS(SELECT * FROM user WHERE username = ?)", username).Scan(&exists)
	return exists
}
