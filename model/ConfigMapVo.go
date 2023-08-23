package model

import (
	"devopscenter/helper"
)

type ConfigBase struct {
	Env       string `json:"env"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type ConfigMapJson struct {
	ConfigBase
	ConfigMapData
}

type ConfigMapData struct {
	Data []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"data"`
}

type ConfigMap struct {
	Id int `json:"id"`
	ConfigBase
	Data string `json:"data"`
}

func (c ConfigMap) Insert(base ConfigBase, data string) (bool, error) {
	query := "INSERT INTO configmap (env, namespace, name, data) VALUES (?, ?, ?, ?)"
	mysqlEngine := helper.SqlContext
	result, err := mysqlEngine.Exec(query, base.Env, base.Namespace, base.Name, data)
	if err != nil {
		return false, nil
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

func (c ConfigMap) List(env, namespace string) ([]*ConfigMap, error) {
	query := "SELECT id, name, data FROM configmap WHERE env=? AND namespace=?"
	mysqlEngine := helper.SqlContext
	rows, err := mysqlEngine.Query(query, env, namespace)
	if err != nil {
		return nil, err
	}
	data := make([]*ConfigMap, 0)
	for rows.Next() {
		obj := &ConfigMap{}
		err := rows.Scan(&obj.Id, &obj.Name, &obj.Data)
		if err != nil {
			return nil, err
		}
		data = append(data, obj)
	}
	return data, nil
}
