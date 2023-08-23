package model

import "devopscenter/helper"

type IngressBase struct {
	Env       string `json:"env"`
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Domain    string `json:"domain"`
}

type IngressCreate struct {
	IngressBase
	Rules []struct {
		Path          string `json:"path"`
		TargetPort    int    `json:"target_port"`
		TargetService string `json:"target_service"`
	} `json:"rules"`
}

type Ingress struct {
	Id int `json:"id"`
	IngressBase
	Data string `json:"data"`
}

func (i Ingress) Insert(env, namespace, name, domain, data string) (bool, error) {
	query := "INSERT INTO ingress(env, namespace, name, domain, data) VALUES (?, ?, ?, ?, ?)"
	mysqlEngine := helper.SqlContext
	result, err := mysqlEngine.Exec(query, env, namespace, name, domain, data)
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

func (i Ingress) List(env, namespace string) ([]*Ingress, error) {
	query := "SELECT id, name, domain, data FROM ingress WHERE env=? AND namespace=?"
	mysqlEngine := helper.SqlContext
	rows, err := mysqlEngine.Query(query, env, namespace)
	if err != nil {
		return nil, err
	}
	data := make([]*Ingress, 0)
	for rows.Next() {
		obj := &Ingress{}
		err := rows.Scan(&obj.Id, &obj.Name, &obj.Domain, &obj.Data)
		if err != nil {
			return nil, err
		}
		data = append(data, obj)
	}
	return data, nil
}
