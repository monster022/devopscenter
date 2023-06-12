package model

import (
	"database/sql"
	"devopscenter/helper"
	"log"
	"time"
)

type Application struct {
	Name        string `json:"name"`
	Language    string `json:"language"`
	BuildPath   string `json:"build_path"`
	PackageName string `json:"package_name"`
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
	Language      string `json:"language" db:"language"`
	BuildPath     string `json:"build_path" db:"build_path"`
	PackageName   string `json:"package_name" db:"package_name"`
}

type ProjectDetail struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Params  string `json:"params"`
	Project string `json:"project"`
	Time    string `json:"time"`
	Message string `json:"message"`
	JobName string `json:"job_name"`
	JobId   int    `json:"job_id"`
}

func (p *Project) Insert() bool {
	mysqlEngine := helper.SqlContext
	var count int
	mysqlEngine.QueryRow("select count(*) from project where project_name = ?", p.ProjectName).Scan(&count)
	if count != 0 {
		return false
	}
	_, err := mysqlEngine.Exec("insert into project(project_id, project_name, project_repo, project_status, project_number, language, build_path, package_name) values (?, ?, ?, ?, ?, ?, ?, ?)",
		p.ProjectId, p.ProjectName, p.ProjectRepo, p.ProjectStatus, p.ProjectNumber, p.Language, p.BuildPath, p.PackageName)
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
		err = rows.Scan(&obj.Id, &obj.ProjectId, &obj.ProjectName, &obj.ProjectRepo, &obj.ProjectStatus, &obj.ProjectNumber, &obj.Language, &obj.BuildPath, &obj.PackageName)
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

func (p *Project) Edit(name, buildPath, packageName string) bool {
	mysqlEngine := helper.SqlContext
	_, err := mysqlEngine.Exec("UPDATE project SET build_path=?, package_name=? WHERE project_name=?", buildPath, packageName, name)
	if err != nil {
		return false
	}
	return true
}

func (d ProjectDetail) List(project string, page, size int) (data []*ProjectDetail) {
	mysqlEngine := helper.SqlContext
	rows, err := mysqlEngine.Query("SELECT id, name, job_name, job_id, params, project, message, time FROM build_info WHERE project = ? ORDER BY id DESC limit ? offset ?", project, size, (page-1)*size)
	if err == sql.ErrNoRows {
		log.Printf("Non Rows")
	}
	for rows.Next() {
		obj := &ProjectDetail{}
		err = rows.Scan(&obj.Id, &obj.Name, &obj.JobName, &obj.JobId, &obj.Params, &obj.Project, &obj.Message, &obj.Time)
		if err != nil {
			log.Fatalln(err)
		}
		data = append(data, obj)
	}
	defer rows.Close()
	return data
}

func (d ProjectDetail) Create() bool {
	mysqlEngine := helper.SqlContext
	_, err := mysqlEngine.Exec("INSERT INTO build_info(project, name, job_name, job_id, message, params) VALUES (?, ?, ?, ?, ?, ?)", d.Project, d.Name, d.JobName, d.JobId, d.Message, d.Params)
	if err != nil {
		return false
	}
	return true
}
