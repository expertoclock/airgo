# 🚀 AirGo: High-Performance LAN File Sharing

AirGo is a lightweight, futuristic, and secure file-sharing tool designed specifically for local area networks (LAN). Built with **Go**, **Gin**, and **Tailwind CSS**, it combines extreme backend performance with a premium user experience.

![Go Version](https://img.shields.io/badge/Go-1.23-00ADD8?style=flat-square&logo=go)
![Docker](https://img.shields.io/badge/Docker-Enabled-2496ED?style=flat-square&logo=docker)
![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)

---

## ✨ Features

- **Blazing Fast:** Powered by Go's high-concurrency model.
- **Futuristic UI:** Glassmorphism design with Dark Mode, powered by Alpine.js and Tailwind CSS.
- **Drag & Drop:** Effortless file uploads with live progress bars.
- **DevOps Ready:** 
  - Multi-stage Docker builds.
  - CI/CD via GitHub Actions.
  - Resource-limited container orchestration.
  - Graceful shutdown and health checks.
- **Security First:** Runs as a non-root user with dropped kernel capabilities.

---

## 🛠 Tech Stack

- **Backend:** [Go](https://go.dev/) with [Gin Web Framework](https://gin-gonic.com/)
- **Frontend:** [Alpine.js](https://alpinejs.dev/), [Tailwind CSS](https://tailwindcss.com/), [Axios](https://axios-http.com/)
- **DevOps:** Docker, Docker Compose, GitHub Actions, Make
- **Icons:** FontAwesome 6

---

## 🚀 Getting Started & Workflows

AirGo supports two **separate, alternative** workflows depending on your needs. You should run *either* the Docker workflow *or* the Kubernetes workflow, not both at the same time. 

### 1. The Developer Workflow (Local Docker)
*Best for quick testing, development, and simple local file sharing.*

1. **Start the application:**
   ```bash
   make up
   ```
2. **Access the tool:** Open `http://localhost:8081` (or your LAN IP).
3. **Stop & Clean:**
   ```bash
   make down
   ```

### 2. The Administrator Workflow (Kubernetes & Monitoring)
*Best for studying DevOps, CI/CD, and running a resilient cluster with monitoring.*

1. **Deploy to Minikube:**
   ```bash
   make k8s-deploy
   ```
   *(This starts Minikube, builds the image inside the cluster, and applies replicas).*
2. **Access the tool:**
   ```bash
   make k8s-url
   ```
   *(This will automatically provide the Kubernetes exposed URL and open it in your browser).*
3. **Verify & Monitor:**
   - Check status: `make k8s-status`
   - Open Grafana monitoring: `make k8s-grafana`
3. **Teardown Cluster App:**
   ```bash
   make k8s-down
   ```

---

## 📱 LAN Access (Mobile)

To share files with your mobile device:
1. Ensure both devices are on the same Wi-Fi.
2. Find your laptop's IP address: `hostname -I`.
3. Type `http://<IP_ADDRESS>:8081` into your phone's browser.
4. **Firewall Note:** If you can't connect, you may need to allow port 8081:
   `sudo ufw allow 8081/tcp`


---

## 📂 Project Structure

```text
.
├── .github/workflows/ # CI/CD Pipeline
├── templates/         # Frontend HTML/JS/CSS
├── uploads/           # Persistent storage for files
├── main.go            # Backend entry point
├── Dockerfile         # Multi-stage build definition
├── compose.yml        # Orchestration & Volume mapping
├── Makefile           # Automation shortcuts
└── go.mod             # Go module definition
```

---

## ⚙️ Configuration

AirGo can be configured using environment variables in the `compose.yml` or a `.env` file:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Internal port the app listens on | `8080` |
| `MAX_UPLOAD_SIZE_MB` | Maximum allowed file size | `500` |
| `GIN_MODE` | Gin framework mode (`debug` or `release`) | `release` |

---

## 🏗 DevOps Commands

| Command | Action |
|---------|--------|
| `make up` | Build and start the containerized app |
| `make down` | Stop and remove containers |
| `make logs` | Follow application logs |
| `make clean` | Remove binaries and clear uploads |
| `make prune` | Remove all unused Docker build cache |
| `make test` | Run backend unit tests |
| `make k8s-deploy` | Build and deploy the app to Minikube Kubernetes |
| `make k8s-url` | Get the URL to access AirGo exposed by Minikube |
| `make k8s-down` | Remove the app from Kubernetes |
| `make k8s-status` | Check Kubernetes pods and services |
| `make k8s-grafana` | Open the Grafana monitoring dashboard |

---

## 🔒 Security

- **Non-Root Execution:** The Docker container runs as `appuser` (UID 1000).
- **Graceful Shutdown:** Handles `SIGTERM` to allow finishing active uploads.
- **Resource Constraints:** Limits defined in `compose.yml` to prevent DoS.
- **Static Analysis:** `go vet` is integrated into the CI pipeline.

---

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.

---
*Built for learning and high-performance sharing.*
