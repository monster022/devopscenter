package model

import "devopscenter/helper"

type PubMachine struct {
	Id          int    `json:"id"`
	MachineId   int    `json:"machine_id"`
	MachineType string `json:"machine_type"`
	ProjectId   int    `json:"project_id"`
}

func (p PubMachine) List(id int) ([]string, error) {
	query := "SELECT instance_ip FROM publish_machine a INNER JOIN machine m on m.id = a.machine_id INNER JOIN project p on p.id = a.project_id WHERE a.project_id = ?"
	mysqlEngine := helper.SqlContext
	rows, err := mysqlEngine.Query(query, id)
	if err != nil {
		return nil, err
	}
	data := make([]string, 0)
	for rows.Next() {
		var obj string
		if ok := rows.Scan(&obj); ok != nil {
			return nil, ok
		}
		data = append(data, obj)
	}
	return data, nil
}
