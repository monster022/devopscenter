package model

import (
	"database/sql"
	"devopscenter/helper"
	"log"
)

type DeploymentAdd struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Replicas  int    `json:"replicas"`
	Image     string `json:"image"`
	Env       string `json:"env"`
}

type DeploymentImage struct {
	ImageSource    string `json:"image_source"`
	ContainerName  string `json:"container_name"`
	DeploymentName string `json:"deployment_name"`
	Namespace      string `json:"namespace"`
	Env            string `json:"env"`
	PublishType    string `json:"publish_type"`
	CreateBy       string `json:"create_by"`
}

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
