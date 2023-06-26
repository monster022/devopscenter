package model

import (
	"database/sql"
	"devopscenter/helper"
	"log"
)

/*
CREATE TABLE `machine`  (
	`id` int NOT NULL AUTO_INCREMENT COMMENT 'ID',
	`instance_name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '备注信息',
	`instance_ip` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '机器的IP地址',
	`instance_username` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '机器的用户名',
	`instance_password` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '机器的密码',
	`instance_cpu` int NULL DEFAULT NULL COMMENT '机器的cpu数量',
	`instance_memory` int NULL DEFAULT NULL COMMENT '机器的内存数量',
	`instance_tag` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '机器的类型',
	PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;
*/

type Machine struct {
	Id               int    `json:"id" db:"id"`
	InstanceName     string `json:"instance_name" db:"instance_name"`
	InstanceIp       string `json:"instance_ip" db:"instance_ip"`
	InstanceUsername string `json:"instance_username" db:"instance_username"`
	InstancePassword string `json:"instance_password" db:"instance_password"`
	InstanceCpu      int    `json:"instance_cpu" db:"instance_cpu"`
	InstanceMemory   int    `json:"instance_memory" db:"instance_memory"`
	InstanceTag      string `json:"instance_tag" db:"instance_tag"`
}

func (m *Machine) List(page int, size int) (data []*Machine) {
	mysqlEngine := helper.SqlContext
	data = make([]*Machine, 0)
	rows, err := mysqlEngine.Query("select id, instance_name, instance_ip, instance_username, instance_cpu, instance_memory, instance_tag from machine limit ? offset ?", size, (page-1)*size)
	if err == sql.ErrNoRows {
		log.Printf("Non Rows")
	}
	for rows.Next() {
		obj := &Machine{}
		err = rows.Scan(&obj.Id, &obj.InstanceName, &obj.InstanceIp, &obj.InstanceUsername, &obj.InstanceCpu, &obj.InstanceMemory, &obj.InstanceTag)
		if err != nil {
			log.Fatalln(err)
		}
		data = append(data, obj)
	}
	defer rows.Close()
	return data
}

func (m *Machine) Total() (total int) {
	mysqlEngine := helper.SqlContext
	mysqlEngine.QueryRow("select count(*) from machine").Scan(&total)
	return total
}

func (m *Machine) Insert() bool {
	mysqlEngine := helper.SqlContext
	_, err := mysqlEngine.Exec("insert into machine(instance_name, instance_ip, instance_username, instance_password, instance_cpu, instance_memory, instance_tag) values(?, ?, ?, ?, ?, ?, ?)",
		m.InstanceName, m.InstanceIp, m.InstanceUsername, m.InstancePassword, m.InstanceCpu, m.InstanceMemory, m.InstanceTag)
	if err != nil {
		return false
	}
	return true
}

func (m *Machine) Delete(id int) bool {
	mysqlEngine := helper.SqlContext
	_, err := mysqlEngine.Exec("delete from machine where id = ?", id)
	if err != nil {
		return false
	}
	return true
}

func (m *Machine) Update() (result sql.Result, err error) {
	mysqlEngine := helper.SqlContext
	result, err = mysqlEngine.Exec("update machine set instance_name = ?, instance_ip = ?, instance_username= ?, instance_password = ? , instance_cpu = ?, instance_memory = ? where id = ?", m.InstanceName, m.InstanceIp, m.InstanceUsername, m.InstancePassword, m.InstanceCpu, m.InstanceMemory, m.Id)
	return result, err
}

func (m Machine) PatchName(id int, name string) (result sql.Result, err error) {
	mysqlEngine := helper.SqlContext
	result, err = mysqlEngine.Exec("UPDATE machine SET instance_name = ? WHERE id = ?", name, id)
	return result, err
}

func (m *Machine) PasswordList(id int) (p string) {
	mysqlEngine := helper.SqlContext
	mysqlEngine.QueryRow("select instance_password from machine where id = ?", id).Scan(&p)
	return p
}

func (m *Machine) PasswordByIp(ip string) (string, error) {
	var password string
	query := "SELECT instance_password FROM machine WHERE instance_ip = ?"
	err := helper.SqlContext.QueryRow(query, ip).Scan(&password)
	if err != nil {
		return "", err
	}
	return password, nil
}
