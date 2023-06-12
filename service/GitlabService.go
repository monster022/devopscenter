package service

import (
	"devopscenter/helper"
	"devopscenter/model"
	"github.com/xanzy/go-gitlab"
)

func SearchName(name string) (bool, int, string, error) {
	gitEngine := helper.GitConnect
	project, _, err := gitEngine.Search.Projects(name, &gitlab.SearchOptions{})
	if len(project) == 0 || err != nil {
		return false, 0, "none", err
	}
	for i := 0; i < len(project); i++ {
		if len(project[i].Name) == len(name) {
			//id := project[i].ID
			return true, project[i].ID, project[i].SSHURLToRepo, err
		}
	}
	return false, 0, "none", nil
}

func SearchAll(name string) ([]*gitlab.Project, error) {
	gitEngine := helper.GitConnect
	project, _, err := gitEngine.Search.Projects(name, &gitlab.SearchOptions{})
	if err != nil {
		return nil, err
	}
	return project, err
}

func BranchList(id int) ([]*model.Branch, error) {
	gitEngine := helper.GitConnect
	// 分支数量
	options := gitlab.ListOptions{PerPage: 100}
	branch, _, err := gitEngine.Branches.ListBranches(id, &gitlab.ListBranchesOptions{
		ListOptions: options,
	})
	if err != nil {
		return nil, err
	}
	data := make([]*model.Branch, 0)
	for i := 0; i < len(branch); i++ {
		obj := &model.Branch{}
		obj.ShortID = branch[i].Commit.ShortID
		obj.Name = branch[i].Name
		obj.Message = branch[i].Commit.Message
		obj.CommitterName = branch[i].Commit.CommitterName
		obj.CommittedDate = branch[i].Commit.CommittedDate
		data = append(data, obj)
	}
	return data, err
}

func RecordBuildInfo(params *model.JenkinsTemplate, marshalData string, jobId int) bool {

	data := model.ProjectDetail{}
	data.Project = params.Project
	data.Name = params.CreateBy
	data.Params = marshalData
	data.JobName = params.Language + "_Template"
	data.JobId = jobId + 1
	data.Message = "ing"
	if !data.Create() {
		return false
	}
	return true
}
