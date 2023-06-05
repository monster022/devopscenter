package model

import (
	"database/sql"
	"devopscenter/helper"
	"log"
)

/*
CREATE TABLE `ingress`  (
	`id` int NOT NULL AUTO_INCREMENT COMMENT 'ID',
	`name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'ingress名称',
	`hosts` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'ingress域名',
	`paths` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'ingress匹配路径',
	`env` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'ingress对应的环境',
	`namespace` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'ingress的名称空间',
	PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;
*/

type Ingress struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Host      string `json:"host"`
	Path      string `json:"path"`
	Env       string `json:"env"`
	Namespace string `json:"namespace"`
}

func (i *Ingress) List(env, namespace string, page, size int) (data []*Ingress) {
	mysqlEngine := helper.SqlContext
	data = make([]*Ingress, 0)
	rows, err := mysqlEngine.Query("SELECT id, name, host, path, env, namespace FROM ingress WHERE env = ? AND namespace = ? limit ? offset ?", env, namespace, size, (page-1)*size)
	if err == sql.ErrNoRows {
		log.Printf("Non Rows")
	}
	for rows.Next() {
		obj := &Ingress{}
		err = rows.Scan(&obj.Id, &obj.Name, &obj.Host, &obj.Path, &obj.Env, &obj.Namespace)
		if err != nil {
			log.Fatalln(err)
		}
		data = append(data, obj)
	}
	defer rows.Close()
	return data
}

func (i *Ingress) Count(namespace string) (total int) {
	mysqlEngine := helper.SqlContext
	mysqlEngine.QueryRow("SELECT count(*) FROM ingress WHERE namespace = ?", namespace).Scan(&total)
	return total
}

func (i Ingress) Delete(id int) bool {
	mysqlEngine := helper.SqlContext
	if _, err := mysqlEngine.Exec("DELETE FROM ingress WHERE id = ?", id); err != nil {
		return false
	}
	return true
}
