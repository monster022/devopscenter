package model

import (
	"database/sql"
	"devopscenter/helper"
	"log"
)

//type DeploymentAdd struct {
//	Name      string `json:"name"`
//	Namespace string `json:"namespace"`
//	Replicas  int    `json:"replicas"`
//	Image     string `json:"image"`
//	Env       string `json:"env"`
//}

type DeploymentBase struct {
	Env       string `json:"env"`
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

type DeployAdd struct {
	DeploymentBase
	Replicas   int    `json:"replicas"`
	Image      string `json:"image"`
	ConfigName string `json:"configName"`
	Uri        string `json:"uri"`
	Port       int32  `json:"port"`
	Cpu        string `json:"cpu"`
	Mem        string `json:"mem"`
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

type DeploymentV2 struct {
	Id int `json:"id"`
	DeployAdd
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

func (d DeployAdd) Insert(add *DeployAdd) (bool, error) {
	query := "INSERT INTO deploy (env, namespace, name, replicas, image, configName, uri, port, cpu, mem) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	mysqlEngine := helper.SqlContext
	result, err := mysqlEngine.Exec(query, add.Env, add.Namespace, add.Name, add.Replicas, add.Image, add.ConfigName, add.Uri, add.Port, add.Cpu, add.Mem)
	if err != nil {
		return false, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if rowsAffected > 0 {
		return true, nil
	}
	return false, nil
}

func (d DeployAdd) List(env, namespace string) ([]*DeploymentV2, error) {
	query := "SELECT id, name, replicas, image, configName, uri, port, cpu, mem FROM deploy WHERE env=? AND namespace=?"
	mysqlengine := helper.SqlContext
	rows, err := mysqlengine.Query(query, env, namespace)
	if err != nil {
		return nil, err
	}
	data := make([]*DeploymentV2, 0)
	for rows.Next() {
		obj := &DeploymentV2{}
		err := rows.Scan(&obj.Id, &obj.Name, &obj.Replicas, &obj.Image, &obj.ConfigName, &obj.Uri, &obj.Port, &obj.Cpu, &obj.Mem)
		if err != nil {
			return nil, err
		}
		data = append(data, obj)
	}
	return data, nil
}
