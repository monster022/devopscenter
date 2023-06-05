package model

type BuildParams struct {
	Repository  string `json:"repository"`
	Project     string `json:"project"`
	Branch      string `json:"branch"`
	PackageName string `json:"package_name"`
}

type TemplateDotnet struct {
	Env                 string `json:"env"`
	Repository          string `json:"repository"`
	DependentRepository string `json:"dependent_repository"`
	Project             string `json:"project"`
	DependentProject    string `json:"dependent_project"`
	SubName             string `json:"sub_name"`
	Branch              string `json:"branch"`
	DependentBranch     string `json:"dependent_branch"`
	BuildPath           string `json:"build_path"`
	PackageName         string `json:"package_name"`
	ImageSource         string `json:"image_source"`
}
