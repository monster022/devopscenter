package service

import (
	"devopscenter/helper"
	"devopscenter/model"
)

func CheckJob(name string) bool {
	jenkinsEngine := helper.JkConnect
	_, err := jenkinsEngine.GetJob(helper.Ctx, name)
	if err != nil {
		return false
	}
	return true
}

func BuildJob(name string, data *model.JenkinsTemplate) (int64, error) {
	jenkinsEngine := helper.JkConnect
	params := map[string]string{
		"Repository":           data.Repository,
		"Dependent_Repository": data.DependentRepository,
		"Project":              data.Project,
		"Dependent_Project":    data.DependentProject,
		"Sub_Name":             data.SubName,
		"Branch":               data.Branch,
		"ShortID":              data.ShortID,
		"Dependent_Branch":     data.DependentBranch,
		"Build_Path":           data.BuildPath,
		"Package_Name":         data.PackageName,
		"Environment_Unique":   data.Env,
		"Image_Source":         data.ImageSource,
		"AliasName":            data.AliasName,
		"Create_By":            data.CreateBy,
	}
	result, err := jenkinsEngine.BuildJob(helper.Ctx, name, params)
	return result, err
}

func IdJob(name string) (int64, error) {
	jenkinsEngine := helper.JkConnect
	number, err := jenkinsEngine.GetAllBuildIds(helper.Ctx, name)
	return number[0].Number, err
}
