package model

import "devopscenter/helper"

type CronJobBase struct {
	Env       string `json:"env"`
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Image     string `json:"image"`
	Schedule  string `json:"schedule"`
}

type CronJob struct {
	Id int `json:"id"`
	CronJobBase
	Data string `json:"data"`
}

type CronJobCreate struct {
	CronJobBase
	Data []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"data"`
}

func (c CronJob) Insert(env, namespace, name, schedule, image, data string) (bool, error) {
	query := "INSERT INTO cronjob(env, namespace, name, schedule, image, data) VALUES (?, ?, ?, ?, ?, ?)"
	mysqlEngine := helper.SqlContext
	result, err := mysqlEngine.Exec(query, env, namespace, name, schedule, image, data)
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

func (c CronJob) List(env, namespace string) ([]*CronJob, error) {
	query := "SELECT id, name, schedule, image, data FROM cronjob WHERE env=? AND namespace=?"
	mysqlEngine := helper.SqlContext
	rows, err := mysqlEngine.Query(query, env, namespace)
	if err != nil {
		return nil, err
	}
	data := make([]*CronJob, 0)
	for rows.Next() {
		obj := &CronJob{}
		err := rows.Scan(&obj.Id, &obj.Name, &obj.Schedule, &obj.Image, &obj.Data)
		if err != nil {
			return nil, err
		}
		data = append(data, obj)
	}
	return data, nil
}
