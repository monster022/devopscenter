package service

import (
	"bytes"
	"devopscenter/configuration"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"
)

func ImagePull(machine, fromImage, tag string, wgImagePull *sync.WaitGroup) {
	defer wgImagePull.Done()
	URL := "http://" + machine + ":2375/images/create?fromImage=" + url.QueryEscape(fromImage) + "&tag=" + tag
	req, err := http.NewRequest("POST", URL, nil)
	println(configuration.Configs.DockerPort)
	if err != nil {
		println("ImagePull 函数 NewRequest的时候错了")
	}
	req.Header.Set("Content-Type", "application/json")
	Client := &http.Client{}
	res, err := Client.Do(req)
	if err != nil {
		println("ImagePull 函数发送请求失败")
	}
	defer res.Body.Close()
	data, _ := ioutil.ReadAll(res.Body)
	println("ImagePull 返回信息：" + string(data))
	println("ImagePull 函数发送http请求下载镜像, 状态码为: " + res.Status)
}

func ContainerDelete(machine, containerName string, wgContainerDelete *sync.WaitGroup) {
	defer wgContainerDelete.Done()
	req, err := http.NewRequest("DELETE", "http://"+machine+":2375/containers/"+containerName+"?force=1", nil)
	if err != nil {
		println("workerDelete 函数 NewRequest的时候错了")
	}
	Client := &http.Client{}
	res, err := Client.Do(req)
	if err != nil {
		println("ContainerDelete 发送请求失败")
	}
	defer res.Body.Close()
	data, _ := ioutil.ReadAll(res.Body)
	println("ContainerDelete 返回信息：" + string(data))
	println("ContainerDelete 函数发送http请求删除容器, 状态码为: " + res.Status)
}

func ContainerCreate(machine, aspEnvironment, image, containerName, containerPort string, wgContainerCreate *sync.WaitGroup) {
	defer wgContainerCreate.Done()
	jsonBody := []byte(fmt.Sprintf(`{
		"Hostname":"%s",
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
	}`, containerName, aspEnvironment, image, containerName, containerPort))
	req, err := http.NewRequest("POST", "http://"+machine+":2375/containers/create?name="+containerName, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		println("函数workerCreate NewRequest发生错误")
	}
	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		println("ContainerCreate 错误为" + err.Error())
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	println("ContainerCreate 返回信息：" + string(data))
	println("ContainerCreate 函数发送http请求创建容器, 状态码为: " + res.Status)
}

func ContainerStart(machine, containerName string, wgContainerStart *sync.WaitGroup) {
	defer wgContainerStart.Done()
	req, err := http.NewRequest("POST", "http://"+machine+":2375/containers/"+containerName+"/start", nil)
	if err != nil {
		println("ContainerStart 函数 NewRequest的时候错了")
	}
	Client := &http.Client{}
	res, err := Client.Do(req)
	if err != nil {
		println("ContainerCreate 发送请求失败")
	}
	defer res.Body.Close()
	data, _ := ioutil.ReadAll(res.Body)
	println("ContainerStart 返回信息：" + string(data))
	println("ContainerStart 函数发送http请求启动容器, 状态码为: " + res.Status)
}
