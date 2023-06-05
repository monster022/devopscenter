package model

import (
	"database/sql"
	"devopscenter/helper"
	"log"
)

/*
CREATE TABLE `cronjob`  (
	`id` int NOT NULL AUTO_INCREMENT COMMENT 'ID',
	`name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'job名称',
	`schedule` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'job执行时间',
	`image` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'job镜像',
	`env` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'job对应的环境',
	`namespace` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'job的名称空间',
	PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;
*/

type CronJob struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Schedule  string `json:"schedule"`
	Image     string `json:"image"`
	Env       string `json:"env"`
	Namespace string `json:"namespace"`
}

func (c CronJob) List(env, namespace string, page, size int) (data []*CronJob) {
	mysqlEngine := helper.SqlContext
	rows, err := mysqlEngine.Query("select id, name, schedule, image, env, namespace from cronjob where env = ? and namespace = ? limit ? offset ?", env, namespace, size, (page-1)*size)
	if err == sql.ErrNoRows {
		log.Printf("Non Rows")
	}
	for rows.Next() {
		obj := &CronJob{}
		err = rows.Scan(&obj.Id, &obj.Name, &obj.Schedule, &obj.Image, &obj.Env, &obj.Namespace)
		if err != nil {
			log.Fatalln(err)
		}
		data = append(data, obj)
	}
	defer rows.Close()
	return data
}

func (c CronJob) Count(namespace string) (total int) {
	mysqlEngine := helper.SqlContext
	mysqlEngine.QueryRow("select count(*) from cronjob where namespace = ?", namespace).Scan(&total)
	return total
}

func (c CronJob) Delete(id int) bool {
	mysqlEngine := helper.SqlContext
	_, err := mysqlEngine.Exec("delete from cronjob where id = ?", id)
	if err != nil {
		return false
	}
	return true
}
