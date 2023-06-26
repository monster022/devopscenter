package model

import (
	"database/sql"
	"devopscenter/helper"
)

type Operator struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Instance string `json:"instance"`
	Command  string `json:"command"`
	Time     string `json:"time"`
}

func (o Operator) Create(username, instance, command string) (sql.Result, error) {
	mysqlEngine := helper.SqlContext
	stmt, err := mysqlEngine.Prepare("INSERT INTO operator_log (`username`, `instance`, `command`) VALUES (?, ?, ?)")
	if err != nil {
		return nil, err
	}
	result, err := stmt.Exec(username, instance, command)
	if err != nil {
		return nil, err
	}
	return result, nil
}
