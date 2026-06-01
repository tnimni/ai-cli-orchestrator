# Installation Guide: AI CLI Orchestrator Suite

This guide provides step-by-step instructions to install and run the **AI CLI Orchestrator Suite** on a fresh machine.

## 1. Prerequisites
Ensure the following tools are installed on your system.

### macOS (Homebrew)
```bash
brew install docker kind helm python github-cli go
```

### Ubuntu/Debian
```bash
# Docker
sudo apt-get update
sudo apt-get install ca-certificates curl gnupg
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update
sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Kind
[ $(uname -m) = x86_64 ] && curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.22.0/kind-linux-amd64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind

# Helm
curl https://baltocdn.com/helm/signing.asc | gpg --dearmor | sudo tee /usr/share/keyrings/helm.gpg > /dev/null
sudo apt-get install apt-transport-https --yes
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/helm.gpg] https://baltocdn.com/helm/stable/debian/ all main" | sudo tee /etc/apt/sources.list.d/helm-stable-debian.list
sudo apt-get update
sudo apt-get install helm

# Python 3 & GitHub CLI
sudo apt-get install python3 python3-pip gh -y

# Go (Golang)
sudo apt-get install golang -y
```

### Fedora/RHEL/CentOS
```bash
# Docker
sudo dnf -y install dnf-plugins-core
sudo dnf config-manager --add-repo https://download.docker.com/linux/fedora/docker-ce.repo
sudo dnf install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
sudo systemctl start docker

# Kind
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.22.0/kind-linux-amd64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind

# Helm
sudo dnf install helm

# Python 3 & GitHub CLI
sudo dnf install python3 gh -y

# Go (Golang)
sudo dnf install golang -y
```

### Arch Linux
```bash
sudo pacman -S docker kind helm python python-pip github-cli go
sudo systemctl start docker
```

---

## 2. Quick Start: Local Setup with Kind
The easiest way to get started on a fresh machine is using the automated local setup script.

1.  **Clone the Repository:**
    ```bash
    git clone <repository-url>
    cd ai-cli-orchestrator
    ```

2.  **Run the Automated Setup:**
    This script creates a Kind cluster, builds all Docker images, loads them into the cluster, and deploys the Helm chart.
    ```bash
    make kind-setup
    ```

3.  **Start Port-Forwarding:**
    In a **new terminal window**, keep this command running to allow the client to communicate with the manager inside the cluster:
    ```bash
    kubectl port-forward svc/ai-cli-orchestrator 8080:8080
    ```

---

## 3. Manual/Production Deployment
If you are deploying to an existing cluster (EKS, GKE, etc.) or prefer a manual approach:

1.  **Configure your Registry:**
    Set your Docker registry username (default is `ai-cli-suite`):
    ```bash
    export DOCKER_REPO=your-docker-registry-user
    ```

2.  **Build and Push Images:**
    ```bash
    make build-images
    make build-manager
    make push-all
    ```

3.  **Deploy the Helm Chart:**
    ```bash
    make deploy-chart
    ```

---

## 4. Client Tool Configuration
The `client/aicli.py` tool is the primary interface for managing your AI pods.

1.  **Set up a Virtual Environment:**
    ```bash
    cd client
    python3 -m venv venv
    source venv/bin/activate
    pip install requests
    ```

2.  **Configure the Manager URL:**
    Point the client to your port-forwarded service:
    ```bash
    export AICLI_BASE_URL=http://localhost:8080
    ```

---

## 5. Common Operations
You can now manage AI pods using the `make` wrappers in the root directory or the Python client directly.

*   **Deploy a New Pod:**
    ```bash
    # Deploy a Claude pod
    make deploy-pod TYPE=claude NAME=my-claude-pod API_KEY=sk-your-key

    # Deploy a Gemini pod with Vertex AI auth
    make deploy-pod TYPE=gemini NAME=my-gemini-pod USE_VERTEXAI=true

    # Deploy a pod with SSH keys mounted
    make deploy-pod TYPE=claude NAME=my-pod MOUNT_KEYS=true
    ```

*   **Interact with a Pod:**
    This connects your terminal directly to the AI CLI tool running inside the pod.
    ```bash
    make interact NAME=my-claude-pod
    ```

*   **Cleanup:**
    ```bash
    # Delete a specific pod
    make delete-pod NAME=my-claude-pod

    # Tear down the entire local Kind cluster
    make kind-clean
    ```
