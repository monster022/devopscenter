package model

import (
	"database/sql"
	"devopscenter/helper"
	"log"
)

/*
CREATE TABLE `service`  (

		`id` int NOT NULL AUTO_INCREMENT COMMENT 'ID',
		`name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'service名称',
		`port_name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '端口的名称',
	    `port` int NOT NULL DEFAULT 80 COMMENT 'service的端口',
		`target_port` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'service的目标端口',
		`node_port` int NOT NULL DEFAULT 0 COMMENT 'service的端口映射端口',
		`protocol` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'service端口协议',
		`type` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'service的类型',
		`env` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'ingress对应的环境',
		`namespace` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'ingress的名称空间',
		PRIMARY KEY (`id`) USING BTREE

) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;
*/
type ServiceCreate struct {
	Name       string `json:"name"`
	Protocol   string `json:"protocol"`
	Port       int    `json:"port"`
	TargetPort int    `json:"target_port"`
	Type       string `json:"type"`
}

type Service struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	PortName   string `json:"port_name"`
	Port       int    `json:"port"`
	TargetPort string `json:"target_port"`
	NodePort   int    `json:"node_port"`
	Protocol   string `json:"protocol"`
	Type       string `json:"type"`
	Env        string `json:"env"`
	Namespace  string `json:"namespace"`
}

func (s Service) List(env, namespace string, page, size int) (data []*Service) {
	mysqlEngine := helper.SqlContext
	rows, err := mysqlEngine.Query("select id, name, port_name, port, target_port, node_port, protocol, type, env, namespace from service where env = ? and namespace = ? limit ? offset ?", env, namespace, size, (page-1)*size)
	if err == sql.ErrNoRows {
		log.Printf("Non Rows")
	}
	for rows.Next() {
		obj := &Service{}
		err = rows.Scan(&obj.Id, &obj.Name, &obj.PortName, &obj.Port, &obj.TargetPort, &obj.NodePort, &obj.Protocol, &obj.Type, &obj.Env, &obj.Namespace)
		if err != nil {
			log.Fatalln(err)
		}
		data = append(data, obj)
	}
	defer rows.Close()
	return data
}

func (s Service) Count(namespace string) (total int) {
	mysqlEngine := helper.SqlContext
	mysqlEngine.QueryRow("select count(*) from service where namespace = ?", namespace).Scan(&total)
	return total
}

func (s Service) Insert() bool {
	mysqlEngine := helper.SqlContext
	_, err := mysqlEngine.Exec("insert into service(name, port_name, port, target_port, node_port, protocol, type, env, namespace) values (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		s.Name, s.PortName, s.Port, s.TargetPort, s.NodePort, s.Protocol, s.Type, s.Env, s.Namespace)
	if err != nil {
		return false
	}
	return true
}

func (s Service) Delete(id int) bool {
	mysqlEngine := helper.SqlContext
	if _, err := mysqlEngine.Exec("delete from service where id = ?", id); err != nil {
		return false
	}
	return true
}
