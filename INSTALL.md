# Installation Guide: AI CLI Orchestrator Suite

This guide provides step-by-step instructions to install and run the **AI CLI Orchestrator Suite** on a fresh machine.

## 1. Prerequisites
Ensure the following tools are installed on your system:
*   **Docker:** To build and run container images.
*   **Kubernetes Cluster:** [Kind](https://kind.sigs.k8s.io/) is recommended for local setup.
*   **Helm 3:** To deploy the manager service.
*   **Python 3:** To run the client management tool.
*   **Go 1.26+:** (Optional) Only needed if you plan to modify the Manager source code.

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
