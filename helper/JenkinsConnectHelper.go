package helper

import (
	"context"
	"devopscenter/configuration"
	"github.com/bndr/gojenkins"
)

var JkConnect *gojenkins.Jenkins
var Ctx = context.Background()

func JenkinsConnect() {
	jenkins := gojenkins.CreateJenkins(nil, configuration.Configs.JenkinsUrl, configuration.Configs.JenkinsUsername, configuration.Configs.JenkinsPassword)
	_, err := jenkins.Init(Ctx)
	if err != nil {
		panic(err)
	}
	JkConnect = jenkins
}
