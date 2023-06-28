# DevopsCenter

## Introduce
$Automated$ $Publishing$ $Platform$  
![Project](https://github.com/monster022/staticResource/blob/main/Images/ProjectTable.jpg)
![Machine](https://github.com/monster022/staticResource/blob/main/Images/MachineTable.jpg)
![ToDoList](https://github.com/monster022/staticResource/blob/main/Images/ToDoListTable.jpg)
Currently supported projects are `Dotnet` `GoLang` and `Vue`

### Dotnet project structure
```  
[root@localhost demo]# tree -L 2  
.  
└── Demo  
├── Demo.API  
├── Demo.Application  
├── Demo.Common  
├── Demo.Contract  
├── Demo.DataAccess  
├── Demo.Domain  
├── Demo.Entity  
├── Demo.IApplication  
├── Demo.IDomain  
├── Demo.Infrastructure  
├── Demo.IRepository  
├── Demo.Proxy  
├── Demo.Repository  
├── Demo.Services  
├── Demo.sln  
├── Demo.UnitTest  
├── DemoV2.Contract  
├── DemoV2.Flight.Api  
├── Tool  
└── nuget.config  
  
19 directories, 2 files  
```  
### GoLang project structure
```  
[root@localhost devopscenter]# tree -L 1  
.  
├── config.ini  
├── configuration  
├── controller  
├── go.mod  
├── go.sum  
├── helper  
├── main.go  
├── middleware  
├── model  
├── README.md  
├── router  
├── service  
└── utils  
```  
### Vue project structure (vue-admin-template)
```  
[root@localhost devopscenterweb]# tree -L 2  
.  
├── babel.config.js  
├── build  
├── index.html  
├── jest.config.js  
├── jsconfig.json  
├── LICENSE  
├── mock  
├── package.json  
├── package-lock.json  
├── postcss.config.js  
├── public  
├── README.md  
├── src  
├── tests  
└── vue.config.js  
  
5 directories, 10 files  
```  
### Jenkins pipeline Scripts
#### dotnet_Template
```grovy  
def generateVersion() {  
return new Date().format('yyyyMMddHHmmss')  
}  
pipeline {  
agent {  
label 'Test_Linux'  
}  
environment {  
Image_Version = generateVersion()  
Harbor_Url = "harbor.chengdd.cn"  
}  
stages {  
stage('Pull Code') {  
steps {  
sh '''  
git clone -b $Dependent_Branch $Dependent_Repository  
git clone -b $Branch $Repository  
'''  
}  
}  
stage('Dotnet Build') {  
steps {  
sh '''  
cd $Project/$Build_Path  
dotnet restore  
dotnet build  
dotnet publish -c Debug -o out  
  
'''  
}  
}  
stage('Write Dockerfile') {  
steps {  
sh '''  
cd $Project/$Build_Path  
echo \"  
FROM ${Image_Source}  
WORKDIR /opt  
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' >/etc/timezone  
COPY out ./  
ENTRYPOINT [\\"dotnet\\", \\"${Package_Name}.dll\\"] \" > Dockerfile  
'''  
}  
}  
stage('Build & Deploy') {  
steps {  
sh '''  
cd $Project/$Build_Path  
docker build -t ${Harbor_Url}/${Environment_Unique}/${AliasName}${Sub_Name}:$Image_Version ./  
docker push ${Harbor_Url}/${Environment_Unique}/${AliasName}${Sub_Name}:$Image_Version  
'''  
}  
}  
}  
post {  
always {  
deleteDir()  
}  
failure {  
sh '/usr/local/sbin/feishutalk --job-name=${JOB_NAME} --project=${AliasName} --build-display-name=${BUILD_DISPLAY_NAME} --message=失败 --author=${Create_By}'}  
success {  
sh '/usr/local/sbin/feishutalk --job-name=${JOB_NAME} --project=${AliasName} --build-display-name=${BUILD_DISPLAY_NAME} --message=成功 --author=${Create_By}'}  
}  
}  
```  
#### go_Template
```  
def generateVersion() {  
return new Date().format('yyyyMMddHHmmss')  
}  
pipeline {  
agent {  
label 'Test_Linux'  
}  
environment {  
Image_Version = generateVersion()  
Harbor_Url = "harbor.chengdd.cn"  
}  
stages {  
stage('Pull Code') {  
steps {  
sh '''  
git clone -b $Branch $Repository  
'''  
}  
}  
stage('Golang Build') {  
steps {  
sh '''  
cd $Project  
go env -w CGO_ENABLED=0  
go env -w GOPROXY=https://goproxy.cn,direct  
go build -o $Project  
'''  
}  
}  
stage('Write Dockerfile') {  
steps {  
sh '''  
cd $Project  
echo \"  
FROM ${Image_Source}  
WORKDIR /opt  
RUN mkdir logs  
COPY ./${Project} ./  
ENTRYPOINT [\\"./${Project}\\"] \" > Dockerfile  
'''  
}  
}  
stage('Build') {  
steps {  
sh '''  
cd $Project  
docker build -t ${Harbor_Url}/${Environment_Unique}/${Project}${Sub_Name}:$Image_Version ./  
docker push ${Harbor_Url}/${Environment_Unique}/${Project}${Sub_Name}:$Image_Version  
# kubectl --kubeconfig=/root/.kube/${Environment_Unique}config set image deployment/${Project}${Sub_Name} -n ${Environment_Unique} ${Project}${Sub_Name}=${Harbor_Url}/${Environment_Unique}/${Project}${Sub_Name}:$Image_Version  
'''  
}  
}  
}  
post {  
always {  
deleteDir()  
}  
failure {  
sh '/usr/local/sbin/feishutalk --job-name=${JOB_NAME} --project=${AliasName} --build-display-name=${BUILD_DISPLAY_NAME} --message=失败 --author=${Create_By}'}  
success {  
sh '/usr/local/sbin/feishutalk --job-name=${JOB_NAME} --project=${AliasName} --build-display-name=${BUILD_DISPLAY_NAME} --message=成功 --author=${Create_By}'}  
}  
}  
```  
#### vue_Template
```  
def generateVersion() {  
return new Date().format('yyyyMMddHHmmss')  
}  
pipeline {  
agent {  
label 'Test_Linux'  
}  
environment {  
Image_Version = generateVersion()  
Harbor_Url = "harbor.chengdd.cn"  
}  
stages {  
stage('Pull Code') {  
steps {  
sh '''  
git clone -b $Branch $Repository  
'''  
}  
}  
stage('vue Build') {  
steps {  
sh '''  
cd $Project  
npm install  
npm run build:prod  
'''  
}  
}  
stage('Write Dockerfile') {  
steps {  
sh '''  
cd $Project  
echo \"  
FROM ${Image_Source}  
WORKDIR /opt  
ADD ./dist ./ \" > Dockerfile  
'''  
}  
}  
stage('Build') {  
steps {  
sh '''  
cd $Project  
docker build -t ${Harbor_Url}/${Environment_Unique}/${Project}${Sub_Name}:$Image_Version ./  
docker push ${Harbor_Url}/${Environment_Unique}/${Project}${Sub_Name}:$Image_Version  
'''  
}  
}  
}  
post {  
always {  
deleteDir()  
}  
failure {  
sh '/usr/local/sbin/feishutalk --job-name=${JOB_NAME} --project=${AliasName} --build-display-name=${BUILD_DISPLAY_NAME} --message=失败 --author=${Create_By}'}  
success {  
sh '/usr/local/sbin/feishutalk --job-name=${JOB_NAME} --project=${AliasName} --build-display-name=${BUILD_DISPLAY_NAME} --message=成功 --author=${Create_By}'}  
}  
}  
```  

## Installation
### Requirement
+ GoLang Environment, Version 1.20
+ Docker
> Tips: Choose one of the two

### Sourcecode
```shell  
git clone https://github.com/monster022/devopscenter.gitcd devopscentergo build -o devopscenter./devopscenter  
```  

### Docker
```shell  
docker pulldocker run -d -p 8080:8080 --name devopscenter```