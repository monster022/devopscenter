package configuration

import (
	"fmt"
	"gopkg.in/ini.v1"
)

type Configuration struct {
	ServerPort      string `ini:"server_port"`
	MysqlAddress    string `ini:"mysql_address"`
	MysqlUsername   string `ini:"mysql_username"`
	MysqlPassword   string `ini:"mysql_password"`
	MysqlPort       string `ini:"mysql_port"`
	MysqlDatabases  string `ini:"mysql_databases"`
	MysqlCharset    string `ini:"mysql_charset"`
	GitlabUsername  string `ini:"gitlab_username"`
	GitlabPassword  string `ini:"gitlab_password"`
	GitlabUrl       string `ini:"gitlab_url"`
	JenkinsUrl      string `ini:"jenkins_url"`
	JenkinsUsername string `ini:"jenkins_username"`
	JenkinsPassword string `ini:"jenkins_password"`
	HarborUrl       string `ini:"harbor_url"`
	HarborUsername  string `ini:"harbor_username"`
	HarborPassword  string `ini:"harbor_password"`
	LdapPort        string `ini:"ldap_port"`
	LdapHost        string `ini:"ldap_host"`
	FeishuTalk      string `ini:"feishu_talk"`
	DockerPort      string `ini:"docker_port"`
}

var Configs = Configuration{}

func LoadConfig(name string) {
	cfg, err1 := ini.Load(name)
	if err1 != nil {
		fmt.Println("Load Config File Failed: ", err1)
	}
	if err2 := cfg.MapTo(&Configs); err2 != nil {
		fmt.Println("Analyze Failed: ", err2)
	}
}
