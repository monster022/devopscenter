package model

import (
	"database/sql"
	"devopscenter/helper"
)

type Order struct {
	OrderId      int    `json:"orderId"`
	SubmitName   string `json:"submitName"`
	Priority     string `json:"priority"`
	Message      string `json:"message"`
	TackleName   string `json:"tackleName"`
	Status       string `json:"status"`
	RejectReason string `json:"rejectReason"`
	Date         string `json:"date"`
}

type OrderRequestBody struct {
	SubmitName string `json:"submitName"`
	Priority   string `json:"priority"`
	Message    string `json:"message"`
	TackleName string `json:"tackleName"`
}

func (o Order) ListOrder(page, size int) ([]*Order, error) {
	query := "select * from `order` limit ? offset ?"
	mysqlEngine := helper.SqlContext
	rows, err := mysqlEngine.Query(query, size, (page-1)*size)
	if err != nil {
		return nil, err
	}
	data := make([]*Order, 0)
	for rows.Next() {
		var obj = &Order{}
		if ok := rows.Scan(&obj.OrderId, &obj.SubmitName, &obj.Priority, &obj.Message, &obj.TackleName, &obj.Status, &obj.RejectReason, &obj.Date); ok != nil {
			return nil, ok
		}
		data = append(data, obj)
	}
	return data, nil
}

func (o Order) ListTackleName(page, size int, tackleName string) ([]*Order, error) {
	query := "select * from `order` where tackleName = ? limit ? offset ?"
	mysqlEngine := helper.SqlContext
	rows, err := mysqlEngine.Query(query, tackleName, size, (page-1)*size)
	if err != nil {
		return nil, err
	}
	data := make([]*Order, 0)
	for rows.Next() {
		var obj = &Order{}
		if ok := rows.Scan(&obj.OrderId, &obj.SubmitName, &obj.Priority, &obj.Message, &obj.TackleName, &obj.Status, &obj.RejectReason, &obj.Date); ok != nil {
			return nil, ok
		}
		data = append(data, obj)
	}
	return data, nil
}

func (o Order) ListTackleNameCount(tackleName string) (total int, err error) {
	queryDoing := "select count(*) from `order` WHERE tackleName = ? AND status in ('doing', 'await')"
	mysqlEngine := helper.SqlContext
	err = mysqlEngine.QueryRow(queryDoing, tackleName).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (o Order) ListSubmitName(page, size int, submitName string) ([]*Order, error) {
	query := "select * from `order` where submitName = ? limit ? offset ?"
	mysqlEngine := helper.SqlContext
	rows, err := mysqlEngine.Query(query, submitName, size, (page-1)*size)
	if err != nil {
		return nil, err
	}
	data := make([]*Order, 0)
	for rows.Next() {
		var obj = &Order{}
		if ok := rows.Scan(&obj.OrderId, &obj.SubmitName, &obj.Priority, &obj.Message, &obj.TackleName, &obj.Status, &obj.RejectReason, &obj.Date); ok != nil {
			return nil, ok
		}
		data = append(data, obj)
	}
	return data, nil
}

func (o Order) PatchOrderStatus(orderId int, status string) (sql.Result, error) {
	mysqlEngine := helper.SqlContext
	stmt, err := mysqlEngine.Prepare("UPDATE `order` SET status = ? WHERE orderId = ?")
	if err != nil {
		return nil, err
	}
	result, err := stmt.Exec(status, orderId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (o Order) PostOrderStatusReject(orderId int, status, rejectReason string) (sql.Result, error) {
	mysqlEngine := helper.SqlContext
	stmt, err := mysqlEngine.Prepare("UPDATE `order` SET status = ?, rejectReason = ? WHERE orderId = ?")
	if err != nil {
		return nil, err
	}
	result, err := stmt.Exec(status, rejectReason, orderId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (o Order) CreateOrder(data *OrderRequestBody) (sql.Result, error) {
	mysqlEngine := helper.SqlContext
	stmt, err := mysqlEngine.Prepare("INSERT INTO `order` (`submitName`, `priority`, `message`, `tackleName`, `status`, `rejectReason`) VALUES (?, ?, ?, ?, 'await', ' ')")
	if err != nil {
		return nil, err
	}
	result, err := stmt.Exec(data.SubmitName, data.Priority, data.Message, data.TackleName)
	if err != nil {
		return nil, err
	}
	return result, nil
}
