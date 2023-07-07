package service

import (
	"bytes"
	"devopscenter/configuration"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func CreateContainer(machine, Hostname, AspEnvironment, Image, HostPort string, statusChan chan int) {
	// 创建容器，0 表示创建时失败
	ContainerCreateUrl := "http://" + machine + ":" + configuration.Configs.DockerPort + "/containers/create?name=" + Hostname
	jsonBody := []byte(fmt.Sprintf(`{
		"Hostname": "%s",
		"User": "root",
		"Env": [
        	"ASPNETCORE_ENVIRONMENT=%s"
		],
		"ExposedPorts": {
	 		"80/tcp": {}
	 	},
		"WorkingDir": "/opt",
		"Image": "%s",
		"Labels": {
			"app": "%s"
		},
		"HostConfig": {
			"Binds": [
            	"/etc/localtime:/etc/localtime:ro"
        	],
	    	"PortBindings": {
				"80/tcp": [
					{
						"HostIp": "0.0.0.0",
						"HostPort":"%s"
					}
				]
			},
			"Privileged": true
		}
	}`, Hostname, AspEnvironment, Image, AspEnvironment, HostPort))
	ContainerCreateReq, err := http.NewRequest("POST", ContainerCreateUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		statusChan <- 0
		return
	}
	ContainerCreateReq.Header.Set("Content-Type", "application/json")
	ContainerCreateClient := &http.Client{}
	ContainerCreateRes, err := ContainerCreateClient.Do(ContainerCreateReq)
	if ContainerCreateRes.StatusCode != http.StatusCreated || err != nil {
		statusChan <- 0
		return
	}
	resBody, err := ioutil.ReadAll(ContainerCreateRes.Body)
	if err != nil {
		statusChan <- 0
		return
	}
	var data struct {
		Id       string        `json:"Id"`
		Warnings []interface{} `json:"Warnings"`
	}
	if err := json.Unmarshal(resBody, &data); err != nil {
		statusChan <- 0
		return
	}
	statusChan <- 88
}

func StartContainer(machine, containerName string, statusChan chan int) {
	// 启动容器， 1 表示启动时失败
	ContainerStartUrl := "http://" + machine + ":" + configuration.Configs.DockerPort + "/containers/" + containerName + "/start"
	ContainerStartReq, err := http.NewRequest("POST", ContainerStartUrl, nil)
	ContainerStartReq.Header.Set("Content-Type", "application/json")
	if err != nil {
		statusChan <- 1
		return
	}
	ContainerClient := &http.Client{}
	ContainerStartRes, err := ContainerClient.Do(ContainerStartReq)
	if err != nil || ContainerStartRes.StatusCode != http.StatusNoContent {
		statusChan <- 1
		return
	}
	statusChan <- 88
}

func DeleteContainer(machine, containerName string, statusChan chan int) {
	// 删除容器， 2 表示删除容器失败
	ContainerDeleteUrl := "http://" + machine + ":" + configuration.Configs.DockerPort + "/containers/" + containerName + "?force=1"
	ContainerDeleteReq, err := http.NewRequest("DELETE", ContainerDeleteUrl, nil)
	if err != nil {
		statusChan <- 3
		return
	}
	ContainerDeleteClient := &http.Client{}
	ContainerDeleteRes, _ := ContainerDeleteClient.Do(ContainerDeleteReq)

	if ContainerDeleteRes.StatusCode == http.StatusNotFound {
		statusChan <- 88
		return
	}
	if ContainerDeleteRes.StatusCode == http.StatusNoContent {
		statusChan <- 88
		return
	}
	statusChan <- 2
}
