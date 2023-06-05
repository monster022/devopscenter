package helper

import (
	"devopscenter/configuration"
	"github.com/xanzy/go-gitlab"
)

var GitConnect *gitlab.Client

func GitlabConnect() {
	client, err := gitlab.NewBasicAuthClient(configuration.Configs.GitlabUsername, configuration.Configs.GitlabPassword,
		gitlab.WithBaseURL(configuration.Configs.GitlabUrl))
	if err != nil {
		panic(err)
	}
	GitConnect = client
}
