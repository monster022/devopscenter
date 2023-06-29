# DevopsCenter

## Introduce
$Automated$ $Publishing$ $Platform$  
Currently supported projects are `Dotnet` `GoLang` and `Vue`  

### Core Component
+ Jenkins
+ GitLab
+ Kubernetes
+ [FeishuTalk](https://github.com/monster022/feishutalk)

#### Jenkins
##### Version 
$2.332.2$
##### Pipeline
###### dotnet_Template
```groovy
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
           sh '/usr/local/sbin/feishutalk --job-name=${JOB_NAME} --project=${AliasName} --build-display-name=${BUILD_DISPLAY_NAME} --message=失败 --author=${Create_By}'
        }
        success {
            sh '/usr/local/sbin/feishutalk --job-name=${JOB_NAME} --project=${AliasName} --build-display-name=${BUILD_DISPLAY_NAME} --message=成功 --author=${Create_By}'
        }
    }
}
```
###### go_Template
```groovy
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
                    go env -w  GOPROXY=https://goproxy.cn,direct
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
                '''
            }
        }
    }
    post {
        always {
            deleteDir()
        }
        failure {
           sh '/usr/local/sbin/feishutalk --job-name=${JOB_NAME} --project=${AliasName} --build-display-name=${BUILD_DISPLAY_NAME} --message=失败 --author=${Create_By}'
        }
        success {
            sh '/usr/local/sbin/feishutalk --job-name=${JOB_NAME} --project=${AliasName} --build-display-name=${BUILD_DISPLAY_NAME} --message=成功 --author=${Create_By}'
        }
    }
}
```
###### vue_Template
```groovy
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
                    npm run build:${Environment_Unique}
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
           sh '/usr/local/sbin/feishutalk --job-name=${JOB_NAME} --project=${AliasName} --build-display-name=${BUILD_DISPLAY_NAME} --message=失败 --author=${Create_By}'
        }
        success {
            sh '/usr/local/sbin/feishutalk --job-name=${JOB_NAME} --project=${AliasName} --build-display-name=${BUILD_DISPLAY_NAME} --message=成功 --author=${Create_By}'
        }
    }
}
```
> Tips:
> Adding different languages and different pipeline scripts can realize the release of different projects

## Installation
### Requirement
+ GoLang Environment, Version 1.20
+ Docker
> Tips: Choose one of the two

### Manual
```shell
# clone code
git clone https://github.com/monster022/devopscenter.git
# enter directory
cd devopscenter
# compile code
go build -o devopscenter
# run code
./devopscenter  
```  

### Docker
```shell
# download images
docker pull $images
# run container
docker run -d -p 8080:8080 --name devopscenter $images
```