
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

