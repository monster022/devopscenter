package model

import (
	"database/sql"
	"devopscenter/helper"
	"log"
	"time"
)

type Application struct {
	Name     string `json:"name"`
	Language string `json:"language"`
}
type Branch struct {
	ShortID       string     `json:"short_id"`
	CommitterName string     `json:"committer_name"`
	CommittedDate *time.Time `json:"committed_date"`
	Message       string     `json:"message"`
	Name          string     `json:"name"`
}

/*
CREATE TABLE `project`  (
	`id` int NOT NULL AUTO_INCREMENT COMMENT 'ID',
	`project_id` int NULL DEFAULT NULL COMMENT '项目ID',
	`project_name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '项目名称',
	`project_repo` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '项目SSH仓库地址',
	`project_status` int NULL DEFAULT 1 COMMENT '项目状态： 1 表示开启, 0 表示关闭',
	`project_number` int NULL DEFAULT 0 COMMENT '项目发布的副本数',
	PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;
*/

type Project struct {
	Id            int    `json:"id" db:"id"`
	ProjectId     int    `json:"project_id" db:"project_id"`
	ProjectName   string `json:"project_name" db:"project_name"`
	ProjectRepo   string `json:"project_repo" db:"project_repo"`
	ProjectStatus int    `json:"project_status" db:"project_status"`
	ProjectNumber int    `json:"project_number" db:"project_number"`
	Language      string `json:"language"`
}

func (p *Project) Insert() bool {
	mysqlEngine := helper.SqlContext
	var count int
	mysqlEngine.QueryRow("select count(*) from project where project_name = ?", p.ProjectName).Scan(&count)
	if count != 0 {
		return false
	}
	_, err := mysqlEngine.Exec("insert into project(project_id, project_name, project_repo, project_status, project_number, language) values (?, ?, ?, ?, ?, ?)",
		p.ProjectId, p.ProjectName, p.ProjectRepo, p.ProjectStatus, p.ProjectNumber, p.Language)
	if err != nil {
		return false
	}
	return true
}

func (p *Project) Patch(id int, status int) bool {
	mysqlEngine := helper.SqlContext
	_, err := mysqlEngine.Exec("UPDATE project set project_status = ? where id = ?", status, id)

	if err != nil {
		return false
	}
	return true
}

func (p *Project) List(page int, size int) (data []*Project) {
	mysqlEngine := helper.SqlContext
	rows, err := mysqlEngine.Query("select * from project limit ? offset ?", size, (page-1)*size)
	if err == sql.ErrNoRows {
		log.Printf("Non Rows")
	}
	for rows.Next() {
		obj := &Project{}
		err = rows.Scan(&obj.Id, &obj.ProjectId, &obj.ProjectName, &obj.ProjectRepo, &obj.ProjectStatus, &obj.ProjectNumber, &obj.Language)
		if err != nil {
			log.Fatalln(err)
		}
		data = append(data, obj)
	}
	defer rows.Close()
	return data
}

func (p *Project) Count() (total int) {
	mysqlEngine := helper.SqlContext
	mysqlEngine.QueryRow("select count(*) from project").Scan(&total)
	return total
}

func (p *Project) Delete(id int) bool {
	mysqlEngine := helper.SqlContext
	_, err := mysqlEngine.Exec("delete from project where id = ?", id)
	if err != nil {
		return false
	}
	return true
}
