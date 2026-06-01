# AI CLI Orchestrator Suite

This repository contains a complete orchestration system for running AI CLIs (Claude, Codex, Gemini, Bob) on Kubernetes.

## Repository Structure

- **`Makefile`**: The central entry point for managing the entire suite.
- **`docker-images/`**: Dockerfiles and build scripts for the AI CLI tool images.
- **`manager/`**: Go-based Kubernetes manager service that handles pod lifecycles.
- **`chart/`**: Helm chart to deploy the orchestrator and configure RBAC.
- **`client/`**: Python CLI tool for interacting with the manager from your local machine.

## Getting Started

The entire project is managed through the **root `Makefile`**. You can see all available commands by running:

```bash
make help
```

### 1. Configure your Repository
Set your Docker registry user to ensure all images are tagged correctly:
```bash
export DOCKER_REPO=your-docker-registry-user
```

### 2. Build Everything
Run these from the root directory to build all components:
```bash
# Build all AI CLI images
make build-images

# Or build a specific CLI image
make claude
make gemini

# Build the Go Manager image
make build-manager
```

### 3. Push to Registry
```bash
make push-all
```

### 4. Deploy to Kubernetes
```bash
make deploy-chart
```

### 5. Local Testing with Kind
If you don't have a remote registry, you can test locally using **Kind**:

```bash
# Create cluster, build images, load them, and deploy
make kind-setup

# If you make changes to the manager code:
make kind-refresh-manager

# If you make changes to any CLI Dockerfiles:
make kind-refresh-clis
```

### 6. Manage Pods via Makefile
You can now manage pods directly from the root directory. You can also inject API keys for authentication:

```bash
# Deploy a Claude pod with an API key
make deploy-pod TYPE=claude NAME=my-pod API_KEY=your-anthropic-key

# Deploy a Gemini pod with advanced auth options
make deploy-pod TYPE=gemini NAME=my-gemini-pod USE_VERTEXAI=true USE_GCA=true

# Deploy with Service Account JSON
make deploy-pod TYPE=gemini NAME=my-gemini-pod GOOGLE_AUTH=path/to/service-account.json

# Deploy with SSH keys mounted (useful for Git operations inside the pod)
make deploy-pod TYPE=claude NAME=my-pod MOUNT_KEYS=true

# Interact with a pod as if local
make interact NAME=my-pod

# Restart a pod
make restart-pod NAME=my-pod

# Delete a pod
make delete-pod NAME=my-pod
```

### 7. Mounting Local SSH Keys
The orchestrator supports mounting your local SSH keys into the AI CLI pods. This is particularly useful if you need to perform Git operations (like cloning private repos) from within the pod.

When you use the `MOUNT_KEYS=true` flag:
1.  The manager identifies your local `~/.ssh` directory.
2.  It creates a `HostPath` volume pointing to that directory.
3.  It mounts this volume as **read-only** to `/root/.ssh` inside the CLI container.

**Note:** This feature is primarily designed for local development environments (like Kind or Docker Desktop) where the host's file system is accessible to the Kubernetes nodes.

### 8. Interactive Experience
When you run `make interact NAME=my-pod`, the client:
1.  Detects the CLI tool type (Gemini, Claude, etc.) from the pod's labels.
2.  Connects your terminal directly to the CLI tool running inside the container.
3.  Allows you to use the tool with full TTY support (colors, interactive prompts, etc.).

### 9. Debugging Pods
If a pod is not working as expected, you can check the **Manager logs** to see how it was created:

```bash
kubectl logs -l app=ai-cli-manager
```

You can also check the **Pod logs**:
```bash
kubectl logs <pod-name>
```

Since pods use an entrypoint wrapper, they will stay alive even if the command fails, allowing you to jump in:
```bash
kubectl exec -it <pod-name> -- /bin/bash
```

## Requirements
- Docker
- Kubernetes Cluster (e.g., Docker Desktop, Minikube, EKS)
- Helm 3
- Python 3
- Go 1.26+ (for manager development)
