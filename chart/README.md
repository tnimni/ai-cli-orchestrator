# AI CLI Orchestrator Helm Chart

This chart deploys the AI CLI Manager, which handles the lifecycle of AI CLI Pods in the cluster.

## Prerequisites
- Kubernetes 1.19+
- Helm 3.0+

## Installation

### Using the Root Makefile (Recommended)
From the project root directory:
```bash
make deploy-chart DOCKER_REPO=your-repo
```

### Manual Installation (Alternative)
```bash
helm install ai-cli-orchestrator . --set image.repository=your-repo/ai-cli-manager
```

## Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `replicaCount` | Number of manager replicas | `1` |
| `image.repository` | Manager image repository | `ai-cli-manager` |
| `image.tag` | Manager image tag | `latest` |
| `service.type` | Service type | `ClusterIP` |
| `defaultRepo` | Default repository for AI CLI images | `ai-cli-suite` |
| `serviceAccount.create` | Create a ServiceAccount | `true` |
| `serviceAccount.name` | Name of the ServiceAccount | `ai-cli-manager-sa` |

## RBAC
This chart creates a `Role` and `RoleBinding` that allows the manager to manage `Pods` within its own namespace.
