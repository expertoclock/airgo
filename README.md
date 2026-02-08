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

## 🚀 Getting Started

### Prerequisites

- [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)
- [Make](https://www.gnu.org/software/make/) (Optional, but recommended)

### Quick Start (The DevOps Way)

1. **Clone the repository:**
   ```bash
   git clone git@github.com:expertoclock/airgo.git
   cd airgo
   ```

2. **Start the application:**
   ```bash
   make up
   ```

3. **Access the tool:**
   Open your browser and navigate to `http://localhost:8081`.

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
| `make test` | Run backend unit tests |

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
