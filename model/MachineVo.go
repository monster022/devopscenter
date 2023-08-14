package model

import (
	"database/sql"
	"devopscenter/helper"
)

type Machine struct {
	Id               int    `json:"id" db:"id"`
	InstanceName     string `json:"instance_name" db:"instance_name"`
	InstanceIp       string `json:"instance_ip" db:"instance_ip"`
	InstanceUsername string `json:"instance_username" db:"instance_username"`
	InstancePassword string `json:"instance_password" db:"instance_password"`
	InstanceCpu      int    `json:"instance_cpu" db:"instance_cpu"`
	InstanceMemory   int    `json:"instance_memory" db:"instance_memory"`
	InstanceStatus   string `json:"instance_status" db:"instance_status"`
	InstanceTag      string `json:"instance_tag" db:"instance_tag"`
}

func (m Machine) Insert() bool {
	mysqlEngine := helper.SqlContext
	_, err := mysqlEngine.Exec("insert into machine(instance_name, instance_ip, instance_username, instance_password, instance_cpu, instance_memory, instance_tag) values(?, ?, ?, ?, ?, ?, ?)",
		m.InstanceName, m.InstanceIp, m.InstanceUsername, m.InstancePassword, m.InstanceCpu, m.InstanceMemory, m.InstanceTag)
	if err != nil {
		return false
	}
	return true
}

func (m Machine) Delete(id int) bool {
	mysqlEngine := helper.SqlContext
	_, err := mysqlEngine.Exec("delete from machine where id = ?", id)
	if err != nil {
		return false
	}
	return true
}

func (m Machine) Update() (result sql.Result, err error) {
	mysqlEngine := helper.SqlContext
	result, err = mysqlEngine.Exec("update machine set instance_name = ?, instance_ip = ?, instance_username= ?, instance_password = ? , instance_cpu = ?, instance_memory = ? where id = ?", m.InstanceName, m.InstanceIp, m.InstanceUsername, m.InstancePassword, m.InstanceCpu, m.InstanceMemory, m.Id)
	return result, err
}

func (m Machine) PatchName(id int, name string) (result sql.Result, err error) {
	mysqlEngine := helper.SqlContext
	result, err = mysqlEngine.Exec("UPDATE machine SET instance_name = ? WHERE id = ?", name, id)
	return result, err
}

func (m Machine) PasswordList(id int) (p string) {
	mysqlEngine := helper.SqlContext
	mysqlEngine.QueryRow("select instance_password from machine where id = ?", id).Scan(&p)
	return p
}

func (m Machine) PasswordByIp(ip string) (string, error) {
	var password string
	query := "SELECT instance_password FROM machine WHERE instance_ip = ?"
	err := helper.SqlContext.QueryRow(query, ip).Scan(&password)
	if err != nil {
		return "", err
	}
	return password, nil
}

func (m Machine) VagueSearch(ip string, page, size int) ([]*Machine, error) {
	query := "SELECT id, instance_name, instance_ip, instance_username, instance_cpu, instance_memory, instance_tag FROM machine WHERE instance_ip LIKE CONCAT('%', ?, '%') LIMIT ? OFFSET ?"
	mysqlEngine := helper.SqlContext
	rows, err := mysqlEngine.Query(query, ip, size, (page-1)*size)
	if err != nil {
		return nil, err
	}
	data := make([]*Machine, 0)
	for rows.Next() {
		obj := &Machine{}
		err = rows.Scan(&obj.Id, &obj.InstanceName, &obj.InstanceIp, &obj.InstanceUsername, &obj.InstanceCpu, &obj.InstanceMemory, &obj.InstanceTag)
		if err != nil {
			return nil, err
		}
		data = append(data, obj)
	}
	return data, nil
}

func (m Machine) VagueSearchTotal(ip string) (int, error) {
	query := "SELECT count(*) FROM machine WHERE instance_ip LIKE CONCAT('%', ?, '%')"
	mysqlEngine := helper.SqlContext
	rows := mysqlEngine.QueryRow(query, ip)
	var total int
	err := rows.Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (m Machine) DownloadData() ([]*Machine, error) {
	query := "SELECT id, instance_name, instance_ip, instance_username, instance_password, instance_cpu, instance_memory, instance_status, instance_tag FROM machine"
	mysqlEngine := helper.SqlContext
	row, err := mysqlEngine.Query(query)
	if err != nil {
		return nil, err
	}
	data := make([]*Machine, 0)
	for row.Next() {
		obj := &Machine{}
		err := row.Scan(&obj.Id, &obj.InstanceName, &obj.InstanceIp, &obj.InstanceUsername, &obj.InstancePassword, &obj.InstanceCpu, &obj.InstanceMemory, &obj.InstanceStatus, &obj.InstanceTag)
		if err != nil {
			return nil, err
		}
		data = append(data, obj)
	}
	return data, nil
}
