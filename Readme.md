
#### Docker CMD's

1. 
```
docker inspect my-app:dev | grep Architecture
"Architecture": "arm64"
```

2. 
```
docker run --rm my-app:dev uname -m
```
Expected:
```
aarch64
```

3. 
```
docker run --rm my-app:dev file /app/app
```
Expected:
```
ELF 64-bit LSB executable, ARM aarch64
```

4. 
```
docker run --rm -it my-app:dev ls -l /app/app
```
Expected:
```
-rwxr-xr-x  app   ✅ executable
```
If not run:
```
RUN chmod +x /app/app
```

#### To run image

```
docker build -t go-jenkins-ci .
```

```
docker run -p 8080:8080 --name go-jenkins-ci-container go-jenkins-ci
```

Please use the following password to proceed to installation:
[LF]> 
[LF]> e0094ee3b82d457c82a3ec46f83e6010

#### Test
Wheather docker is running on Jenkins container

```
docker exec -it jenkins bash
docker version
```

#### CMD to create directory in home
```
mkdir -p ~/jenkins_home
```

#### Stop Container and Remove
```
docker stop jenkins
```
```
docker rm jenkins
```
#### Dockerfile for Jenkins

```
FROM jenkins/jenkins:lts

USER root

# Install docker CLI
RUN apt-get update && \
    apt-get install -y docker.io && \
    apt-get clean

USER jenkins
```
📌 This installs only the CLI, not Docker daemon.

#### To run jenkin-docker image with docker daemon (maybe), docker.sock access 

##### Step 1: Find docker group ID (inside container)

```
getent group docker
```
You’ll see something like:
```
docker:x:102:
```
👉 102 is the GID

##### Step 2: Restart Jenkins with docker group

```
docker stop jenkins
docker rm jenkins
```
Run again **with group added**:

```
docker run -d \
  --name jenkins \
  -p 8081:8080 \
  -p 50000:50000 \
  -v ~/jenkins_home:/var/jenkins_home \
  -v /var/run/docker.sock:/var/run/docker.sock \
  --group-add 102 \
  jenkins-with-docker
```

⚠️ Replace 102 if your docker group ID is different.

##### From which line jenkins volume is created 

```
 -v ~/jenkins_home:/var/jenkins_home \
```
So even after creating a new image and running a new container 'Create First Admin User' data is Intact.

And this volume will not appear in docker desktop. It's being saved in machine. 

##### From which line docker.sock access is given 

```
docker.sock 
```
#### Run as root user to avoid permission issue
```
docker run -d \
  --name jenkins \
  -p 8081:8080 \
  -p 50000:50000 \
  -v ~/jenkins_home:/var/jenkins_home \
  -v /var/run/docker.sock:/var/run/docker.sock \
  --user root \
  jenkins-with-docker
```

#### Jenkins pipeline script

```
pipeline {
    agent any

    environment {
        IMAGE_NAME = "go-jenkins-ci"
        CONTAINER_NAME = "go-jenkins-ci"
    }

    stages {

        stage('Checkout') {
            steps {
                git branch: 'main',
                    url: 'https://github.com/Fahad-I-Khan/learn-jenkins.git'
            }
        }
        
        stage('Debug Workspace') {
            steps {
                sh '''
                echo "==== LIST ROOT ===="
                ls -la
        
                echo "==== go.mod ===="
                cat go.mod || echo "go.mod NOT FOUND"
        
                echo "==== pwd ===="
                pwd
                '''
            }
        }


        stage('Go Test (Dockerized)') {
            steps {
                sh '''
                docker run --rm \
                  -v "$WORKSPACE":/app \
                  -w /app \
                  golang:1.25.5-alpine3.23 \
                  sh -c "go test ./..."
                '''
            }
        }

        stage('Docker Build') {
            steps {
                sh 'docker build -t $IMAGE_NAME:latest .'
            }
        }

        stage('Run Container') {
            steps {
                sh '''
                docker rm -f $CONTAINER_NAME || true
                docker run -d \
                  -p 8090:8080 \
                  --name $CONTAINER_NAME \
                  $IMAGE_NAME:latest
                '''
            }
        }
    }
}
```

#### Ask git directly 

```
SHORT_SHA=$(git rev-parse --short HEAD)
```

```
stage('Docker Build (Tests enforced)') {
    steps {
        sh '''
        SHORT_SHA=$(git rev-parse --short HEAD)
        echo "Building image with tag: $SHORT_SHA"

        docker build -t go-jenkins-ci:$SHORT_SHA .
        docker tag go-jenkins-ci:$SHORT_SHA go-jenkins-ci:latest
        '''
    }
}
```
This:
- Always works
- Works locally, in Jenkins, in Docker
- Is CI-tool independent

---------------------------------------

### For Git

When git push was not working.
```
 git config --global push.autoSetupRemote true
 ```