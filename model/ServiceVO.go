package model

import (
	"devopscenter/helper"
)

type ServiceCreateV2 struct {
	Env        string `json:"env"`
	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	Port       int    `json:"port"`
	Protocol   string `json:"protocol"`
	TargetPort int    `json:"target_port"`
	Type       string `json:"type"`
	Deployment string `json:"deployment"`
}

type Service struct {
	Id int `json:"id"`
	ServiceCreateV2
}

func (s Service) Insert(v2 *ServiceCreateV2) (bool, error) {
	query := "INSERT INTO service(env, name, namespace, port, protocol, target_port, type, deployment) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	mysqlEngine := helper.SqlContext
	result, err := mysqlEngine.Exec(query, v2.Env, v2.Name, v2.Namespace, v2.Port, v2.Protocol, v2.TargetPort, v2.Type, v2.Deployment)
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

func (s Service) List(env, namespace string) ([]*Service, error) {
	query := "SELECT id, name, port, protocol, target_port, type, deployment FROM service WHERE env=? AND namespace=?"
	mysqlEngine := helper.SqlContext
	rows, err := mysqlEngine.Query(query, env, namespace)
	if err != nil {
		return nil, err
	}
	data := make([]*Service, 0)
	for rows.Next() {
		obj := &Service{}
		err := rows.Scan(&obj.Id, &obj.Name, &obj.Port, &obj.Protocol, &obj.TargetPort, &obj.Type, &obj.Deployment)
		if err != nil {
			return nil, err
		}
		data = append(data, obj)
	}
	return data, nil
}
