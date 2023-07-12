package docker

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	machine := make([]string, 2)
	machine[0] = "10.11.11.124"
	machine[1] = "10.11.11.134"

	var wgPull sync.WaitGroup
	for _, m := range machine {
		wgPull.Add(1)
		//go service.ImagePull(m, "harbor.chengdd.cn/fat/flight", "20230707_131441-7eb9e780", &wgImagePull)
		go workerPull(m, &wgPull)
	}
	wgPull.Wait()

	var wgDelete sync.WaitGroup
	for _, m := range machine {
		wgDelete.Add(1)
		//go service.ContainerDelete(m, "devopscenter", &wgContainerDelete)
		go workerDelete(m, &wgDelete)
	}
	wgDelete.Wait()

	var wgCreate sync.WaitGroup
	for _, m := range machine {
		wgCreate.Add(1)
		//go service.ContainerCreate(m, "fat", "harbor.chengdd.cn/fat/flight:20230707_131441-7eb9e780", "devopscenter", "899", &wgContainerCreate)
		go workerCreate(m, &wgCreate)
	}
	wgCreate.Wait()

	var wgStart sync.WaitGroup
	for _, m := range machine {
		wgStart.Add(1)
		//go service.ContainerStart(m, "devopscenter", &wgContainerStart)
		go workerStart(m, &wgStart)
	}
	wgStart.Wait()
}

func workerPull(machine string, wgPull *sync.WaitGroup) {
	defer wgPull.Done()
	req, err := http.NewRequest("POST", "http://"+machine+":2375/images/create?fromImage=nginx&tag=latest", nil)
	if err != nil {
		println("workerPull函数 NewRequest的时候错了")
	}
	req.Header.Set("Content-Type", "application/json")
	Client := &http.Client{}
	res, err := Client.Do(req)
	if err != nil {
		println("发送请求失败")
	}
	defer res.Body.Close()
	data, _ := ioutil.ReadAll(res.Body)
	println("返回信息：" + string(data))
	println("workerPull函数发送http请求下载镜像, 状态码为: " + res.Status)
}

func workerDelete(machine string, wgDelete *sync.WaitGroup) {
	defer wgDelete.Done()
	req, err := http.NewRequest("DELETE", "http://"+machine+":2375/containers/devopscenter?force=1", nil)
	if err != nil {
		println("workerDelete函数 NewRequest的时候错了")
	}
	ContainerDeleteClient := &http.Client{}
	res, err := ContainerDeleteClient.Do(req)
	if err != nil {
		println("发送请求失败")
	}
	defer res.Body.Close()
	data, _ := ioutil.ReadAll(res.Body)
	println("返回信息：" + string(data))
	println("workerDelete函数发送http请求删除容器, 状态码为: " + res.Status)
}

func workerCreate(machine string, wgCreate *sync.WaitGroup) {
	defer wgCreate.Done()
	ContainerCreateUrl := "http://" + machine + ":2375/containers/create?name=devopscenter"
	jsonBody := []byte("{\n\t\t\"Hostname\": \"devopscenter\",\n\t\t\"User\": \"root\",\n\t\t\"Env\": [\n        \t\"ASPNETCORE_ENVIRONMENT=Production\"\n\t\t],\n\t\t\"ExposedPorts\": {\n\t \t\t\"80/tcp\": {}\n\t \t},\n\t\t\"WorkingDir\": \"/opt\",\n\t\t\"Image\": \"nginx:latest\",\n\t\t\"Labels\": {\n\t\t\t\"app\": \"devopscenter\"\n\t\t},\n\t\t\"HostConfig\": {\n\t\t\t\"Binds\": [\n            \t\"/etc/localtime:/etc/localtime:ro\"\n        \t],\n\t    \t\"PortBindings\": {\n\t\t\t\t\"80/tcp\": [\n\t\t\t\t\t{\n\t\t\t\t\t\t\"HostIp\": \"0.0.0.0\",\n\t\t\t\t\t\t\"HostPort\":\"899\"\n\t\t\t\t\t}\n\t\t\t\t]\n\t\t\t},\n\t\t\t\"Privileged\": true\n\t\t}\n\t}")
	req, err := http.NewRequest("POST", ContainerCreateUrl, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		println("函数workerCreate NewRequest发生错误")
	}
	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		println("错误为" + err.Error())
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	println("返回信息：" + string(data))
	println("workerCreate函数发送http请求创建容器, 状态码为: " + res.Status)
}

func workerStart(machine string, wgStart *sync.WaitGroup) {
	defer wgStart.Done()
	req, err := http.NewRequest("POST", "http://"+machine+":2375/containers/devopscenter/start", nil)
	if err != nil {
		println("workerStart函数 NewRequest的时候错了")
	}
	Client := &http.Client{}
	res, err := Client.Do(req)
	if err != nil {
		println("发送请求失败")
	}
	defer res.Body.Close()
	data, _ := ioutil.ReadAll(res.Body)
	println("返回信息：" + string(data))
	println("workerStart函数发送http请求启动容器, 状态码为: " + res.Status)
}
