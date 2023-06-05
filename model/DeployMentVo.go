package model

import (
	"database/sql"
	"devopscenter/helper"
	"log"
)

type DeploymentImage struct {
	ImageSource    string `json:"image_source"`
	ContainerName  string `json:"container_name"`
	DeploymentName string `json:"deployment_name"`
	Namespace      string `json:"namespace"`
	Env            string `json:"env"`
}

/*
CREATE TABLE `deployment`  (
	`id` int NOT NULL AUTO_INCREMENT COMMENT 'ID',
	`name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'deployment名称',
    `replicate` int NOT NULL DEFAULT 0 COMMENT 'deployment的副本数',
	`container_name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'Pod名称',
	`image` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'Pod镜像',
	`env` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'ingress对应的环境',
	`namespace` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'ingress的名称空间',
	PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;
*/

type Deployment struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Replicate     int    `json:"replicate"`
	ContainerName string `json:"container_name"`
	Image         string `json:"image"`
	Env           string `json:"env"`
	Namespace     string `json:"namespace"`
}

func (m *Deployment) List(env, namespace string, page, size int) (data []*Deployment) {
	mysqlEngine := helper.SqlContext
	rows, err := mysqlEngine.Query("select id, name, replicate, container_name, image, env, namespace from deployment where env = ? and namespace = ? limit ? offset ?", env, namespace, size, (page-1)*size)
	if err == sql.ErrNoRows {
		log.Printf("Non Rows")
	}
	for rows.Next() {
		obj := &Deployment{}
		err = rows.Scan(&obj.Id, &obj.Name, &obj.Replicate, &obj.ContainerName, &obj.Image, &obj.Env, &obj.Namespace)
		if err != nil {
			log.Fatalln(err)
		}
		data = append(data, obj)
	}
	defer rows.Close()
	return data
}

func (m *Deployment) Count(namespace string) (total int) {
	mysqlEngine := helper.SqlContext
	mysqlEngine.QueryRow("select count(*) from deployment where namespace = ?", namespace).Scan(&total)
	return total
}

func (m *Deployment) Delete(id int) bool {
	mysqlEngine := helper.SqlContext
	_, err := mysqlEngine.Exec("delete from deployment where id = ?", id)
	if err != nil {
		return false
	}
	return true
}

func (m *Deployment) ImagePatch(image, env, namespace, name string) int64 {
	mysqlEngine := helper.SqlContext
	result, err := mysqlEngine.Exec("UPDATE deployment SET image=? WHERE env=? AND namespace=? AND name=?", image, env, namespace, name)
	if err != nil {
		return -1 // 数据库执行出错
	}
	affectedRows, err1 := result.RowsAffected()
	if err1 != nil {
		return -1 // 数据库行数影响出错
	}
	return affectedRows
}
