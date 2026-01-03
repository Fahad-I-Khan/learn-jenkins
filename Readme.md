
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

----------
----------
----------
----------
# Header How to set credential in Jenkins
----------
----------
--------------------
----------


Perfect question 👍
I’ll explain this **slowly and clearly**, in **three layers**:

1. **How to add Docker Hub credentials in Jenkins (click-by-click)**
2. **How Jenkins actually uses them in your pipeline**
3. **How Jenkins stores them securely (and why this is safe)**

No assumptions. No magic.

---

# 1️⃣ How to ADD `dockerhub-creds` in Jenkins (Step-by-Step)

### 🔹 Step 1: Open Jenkins Dashboard

Open Jenkins in your browser.

```
http://localhost:8081
```
8080 or port whatever you set. 
---

### 🔹 Step 2: Go to Credentials

Click in this exact order:

```
Dashboard
 → Manage Jenkins
   → Credentials
```

---

### 🔹 Step 3: Choose Global Credentials Store

You’ll see something like:

```
Stores scoped to Jenkins
 └── (global)
```

Click:

```
(global)
```

Then click:

```
Add Credentials
```

---

### 🔹 Step 4: Fill the Credential Form

Use **exactly this** 👇

| Field           | Value                                      |
| --------------- | ------------------------------------------ |
| **Kind**        | `Username with password`                   |
| **Scope**       | `Global`                                   |
| **Username**    | `fahadkhan2105`                            |
| **Password**    | 🔐 Docker Hub **password OR access token** |
| **ID**          | `dockerhub-creds` ✅                        |
| **Description** | Docker Hub credentials                     |

👉 **Important:**
The **ID must be exactly** `dockerhub-creds`
This is how Jenkins finds it from the pipeline.

Click **Save** ✔️

---

# 2️⃣ How It WORKS in Your Jenkins Pipeline

Now let’s connect the dots.

### 🔹 This is the key part in your Jenkinsfile:

```groovy
withCredentials([
    usernamePassword(
        credentialsId: 'dockerhub-creds',
        usernameVariable: 'DOCKER_USER',
        passwordVariable: 'DOCKER_PASS'
    )
]) {
    sh '''
    echo "$DOCKER_PASS" | docker login -u "$DOCKER_USER" --password-stdin
    docker push fahadkhan2105/go-jenkins-ci:latest
    '''
}
```

### What happens internally (step by step):

1. Jenkins sees:

   ```
   credentialsId: 'dockerhub-creds'
   ```
2. Jenkins:

   * Looks into **its credential store**
   * Fetches username + password
3. Jenkins **injects them into memory only** as:

   ```
   DOCKER_USER
   DOCKER_PASS
   ```
4. Inside the `sh` block:

   * Docker logs in using those values
5. After the block finishes:

   * Variables are **destroyed**
   * Credentials are **gone from memory**

🚫 They are **NOT**:

* Printed in logs
* Stored in the workspace
* Saved in Git
* Exposed to other jobs

---

# 3️⃣ How Jenkins Stores Credentials SECURELY

This is the most important part 👇

---

## 🔐 Where Are Credentials Stored?

Inside Jenkins container at:

```
/var/jenkins_home/credentials.xml
```

But ❗ **NOT in plain text**

---

## 🔐 How Are They Encrypted?

Jenkins encrypts credentials using:

### 1️⃣ `master.key`

```
/var/jenkins_home/secrets/master.key
```

### 2️⃣ `hudson.util.Secret`

Jenkins uses its internal crypto system to:

* Encrypt passwords
* Decrypt only at runtime
* Never expose raw values

Even if someone opens `credentials.xml`, they will see **encrypted blobs**, not passwords.

---

## 🔐 Why This Is Secure

| Risk               | Jenkins Protection        |
| ------------------ | ------------------------- |
| Git repo leak      | ❌ creds never in Git      |
| Console logs       | ❌ masked automatically    |
| File system access | 🔐 encrypted              |
| Pipeline access    | 🔒 only via credential ID |
| Accidental echo    | ❌ masked as `****`        |

---

## 🔐 Best Practice (What You’re Already Doing Right)

✅ Use `credentialsId`, not passwords
✅ Use `withCredentials` block
✅ Use Docker access tokens (optional but recommended)
✅ Never hardcode secrets

This is **industry-standard CI/CD security**.

---

# 🔁 Mental Model (Remember This)

> Jenkins credentials are like a **secure vault**
> Your pipeline gets a **temporary key**
> The key expires **as soon as the stage ends**

---

## ✅ What You Should Do Now

1. Add credentials with ID `dockerhub-creds`
2. Re-run the pipeline
3. Watch:

   * Build ✅
   * Push to Docker Hub ✅
4. Check Docker Hub → image with **Git SHA tag**

---

If you want next, I can show:

* 🔐 Docker Hub **Access Token** (better than password)
* 🔄 Rotate credentials safely
* 🔍 How Jenkins masks secrets in logs
* 🚀 Push to AWS ECR (same concept)

Just tell me.

------------
------------
------------
------------
# Header: How to create access token in Docker Hub/hub.docker
------------
------------
------------
------------

Perfect 👍
Using a **Docker Hub access token** is the **right** and **professional** way. I’ll show you **exact click-by-click**, then explain **why it’s better than a password**.

---

# 1️⃣ What is a Docker Hub Access Token (in simple words)

Think of it like this:

* **Password** → full account access 😬
* **Access token** → limited, revocable key 🔐 (best for CI/CD)

Jenkins should **never** use your real password.

---

# 2️⃣ How to CREATE an Access Token on Docker Hub

### 🔹 Step 1: Login to Docker Hub

Go to:

```
https://hub.docker.com
```

Login using:

```
Username: fahadkhan2105
```

---

### 🔹 Step 2: Open Account Settings

Top-right corner → click your **profile icon**
Then click:

```
Account Settings
```

---

### 🔹 Step 3: Go to Security

In the left sidebar, click:

```
Security
```

You’ll see **Access Tokens** section.

---

### 🔹 Step 4: Create New Access Token

Click:

```
New Access Token
```

Fill details:

| Field                 | Value                   |
| --------------------- | ----------------------- |
| **Token Description** | `jenkins-go-jenkins-ci` |
| **Permissions**       | `Read & Write` ✅        |

Click:

```
Generate
```

---

### 🔹 Step 5: COPY the Token (IMPORTANT ⚠️)

Docker Hub will show something like:

```
dckr_pat_xxxxxxxxxxxxxxxxxxxxx
```

📌 **Copy it immediately**

* You will **never** see it again
* If lost → create a new one

---

# 3️⃣ Store This Token in Jenkins (Correct Way)

Now go back to Jenkins:

```
Dashboard
 → Manage Jenkins
   → Credentials
     → (global)
       → Add Credentials
```

Fill like this 👇

| Field           | Value                               |
| --------------- | ----------------------------------- |
| **Kind**        | Username with password              |
| **Username**    | `fahadkhan2105`                     |
| **Password**    | 🔐 **PASTE ACCESS TOKEN HERE**      |
| **ID**          | `dockerhub-creds`                   |
| **Description** | Docker Hub access token for Jenkins |

Click **Save** ✔️

✅ Jenkins treats this token just like a password
✅ Docker CLI supports tokens natively

---

# 4️⃣ How Jenkins Uses the Token (No Code Change Needed)

Your pipeline already works 👇

```groovy
withCredentials([
  usernamePassword(
    credentialsId: 'dockerhub-creds',
    usernameVariable: 'DOCKER_USER',
    passwordVariable: 'DOCKER_PASS'
  )
]) {
  sh '''
    echo "$DOCKER_PASS" | docker login -u "$DOCKER_USER" --password-stdin
    docker push fahadkhan2105/go-jenkins-ci:latest
  '''
}
```

Docker Hub understands:

* `DOCKER_PASS` = access token ✅
* Same login command
* More secure

---

# 5️⃣ Why Access Tokens Are BETTER (Very Important)

| Feature                          | Password | Access Token |
| -------------------------------- | -------- | ------------ |
| Revoke without changing password | ❌        | ✅            |
| Scope limited                    | ❌        | ✅            |
| Safe for CI/CD                   | ❌        | ✅            |
| Audit & rotate                   | ❌        | ✅            |
| Best practice                    | ❌        | ✅            |

👉 **If Jenkins is compromised**, you:

* Revoke token
* Done
* No password reset needed

---

# 6️⃣ How You Know It’s Working

After pipeline runs successfully:

### ✔ Jenkins log shows:

```
docker push fahadkhan2105/go-jenkins-ci:f4b9aa2
docker push fahadkhan2105/go-jenkins-ci:latest
```

### ✔ Docker Hub shows:

```
fahadkhan2105/go-jenkins-ci
  ├── latest
  └── f4b9aa2
```

Each tag = exact Git commit 🚀

---

# 7️⃣ Pro Tip (Industry Standard)

For every CI system:

* Jenkins → token
* GitHub Actions → token
* GitLab CI → token

❌ Never reuse your real password.

---

If you want next:

* 🔁 Rotate tokens safely
* 🔍 See how Jenkins masks secrets in logs
* 🚀 Push same image to AWS ECR
* 🧠 Why tagging with Git SHA is CI gold

Just say 👍

-----------------------
-----------------------
-----------------------
-----------------------
# How to move .pem file from Download to .ssh
-----------------------
-----------------------
-----------------------
-----------------------

Perfect question 👍 — this is **normal Mac behavior**, nothing is wrong.

When AWS downloads the key, macOS **always puts it in `Downloads/` first**.
You **don’t choose the folder at download time** — you move it **after**.

I’ll show you the **exact, safe, step-by-step way**.

---

## Step 1: Confirm the file is in Downloads

Open Terminal and run:

```bash
ls ~/Downloads | grep .pem
```

You should see something like:

```text
jenkins-ec2-key.pem
```

✅ This confirms the key exists.

---

## Step 2: Create `.ssh` directory (if not already there)

Run:

```bash
mkdir -p ~/.ssh
```

* `-p` means “don’t error if it already exists”

---

## Step 3: Move the `.pem` file to `~/.ssh`

Run this **exact command** (replace filename if needed):

```bash
mv ~/Downloads/jenkins-ec2-key.pem ~/.ssh/
```

Now verify:

```bash
ls ~/.ssh
```

You should see:

```text
jenkins-ec2-key.pem
```

---

## Step 4: Fix permissions (VERY IMPORTANT 🚨)

AWS **will reject the key** if permissions are open.

Run:

```bash
chmod 400 ~/.ssh/jenkins-ec2-key.pem
```

Verify permissions:

```bash
ls -l ~/.ssh/jenkins-ec2-key.pem
```

Expected output:

```text
-r--------  1 yourname  staff  jenkins-ec2-key.pem
```

✅ This means SSH will accept it.

---

## Step 5: Why we do this (important concept)

* `~/.ssh` is the **standard secure location**
* SSH tools **expect keys here**
* Prevents accidental deletion
* Permissions are enforced

This is **industry standard**, not just AWS.

---

## Optional (nice DevOps habit)

You can rename it if you want clarity:

```bash
mv ~/.ssh/jenkins-ec2-key.pem ~/.ssh/ec2-jenkins.pem
```

(Just remember the new name)

---

## Final recap

✔ Download → always goes to `Downloads`
✔ You **move it manually** to `~/.ssh`
✔ Set permissions (`chmod 400`)
✔ Now ready for SSH

---

### Next step (tell me once done):

> “Key is in `~/.ssh` and permissions are set”

Then we’ll:
➡️ Launch EC2
➡️ SSH into it
➡️ Install Docker
➡️ Pull image from Docker Hub
➡️ Run your app 🚀

### Jenkins file's Not Good for multi stage 
```
pipeline {
    agent any

    environment {
        IMAGE_NAME = "fahadkhan2105/go-jenkins-ci"
    }

    stages {

        stage('Checkout') {
            steps {
                git branch: 'main',
                    url: 'https://github.com/Fahad-I-Khan/learn-jenkins.git'
            }
        }

        stage('Build Docker Image (Tests enforced)') {
            steps {
                sh '''
                    SHORT_SHA=$(git rev-parse --short HEAD)
                    echo "Git commit SHA: $SHORT_SHA"

                    docker build -t $IMAGE_NAME:$SHORT_SHA .
                    docker tag $IMAGE_NAME:$SHORT_SHA $IMAGE_NAME:latest
                '''
            }
        }

        stage('Push to Docker Hub') {
            steps {
                withCredentials([usernamePassword(
                    credentialsId: 'dockerhub-creds',
                    usernameVariable: 'DOCKER_USER',
                    passwordVariable: 'DOCKER_PASS'
                )]) {
                    sh '''
                        echo "$DOCKER_PASS" | docker login -u "$DOCKER_USER" --password-stdin

                        SHORT_SHA=$(git rev-parse --short HEAD)

                        docker push $IMAGE_NAME:$SHORT_SHA
                        docker push $IMAGE_NAME:latest

                        docker logout
                    '''
                }
            }
        }

        stage('Run Container') {
            steps {
                sh '''
                    SHORT_SHA=$(git rev-parse --short HEAD)

                    docker rm -f go-jenkins-ci || true
                    docker run -d \
                      --name go-jenkins-ci \
                      -p 8080:8080 \
                      $IMAGE_NAME:$SHORT_SHA

                    sleep 3
                    curl -f http://localhost:8080/health
                '''
            }
        }
    }
}
```

### For Git

When git push was not working.
```
 git config --global push.autoSetupRemote true
 ```