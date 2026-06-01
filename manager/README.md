# AI CLI Manager

A Go-based REST API to manage the lifecycle of AI CLI Pods in a Kubernetes cluster.

## API Endpoints

### Create a Pod
- **URL**: `/pods`
- **Method**: `POST`
- **Body**:
  ```json
  {
    "name": "my-claude-pod",
    "type": "claude",
    "repo": "optional-repo-override",
    "apiKey": "optional-api-key-to-inject"
  }
  ```

### Restart a Pod
- **URL**: `/pods/{name}/restart`
- **Method**: `PUT`

### Update a Pod
- **URL**: `/pods/{name}`
- **Method**: `PUT`
- **Body**:
  ```json
  {
    "type": "claude",
    "repo": "new-repo-override"
  }
  ```

### Delete a Pod
- **URL**: `/pods/{name}`
- **Method**: `DELETE`

## Build & Deploy

It is highly recommended to use the **root Makefile** from the project root to manage the manager and its deployment:

### From Project Root (Recommended)
```bash
# Build the manager
make build-manager DOCKER_REPO=your-repo

# Deploy the manager via Helm
make deploy-chart DOCKER_REPO=your-repo
```

### Manual Build (Alternative)
```bash
docker build -t your-repo/ai-cli-manager:latest .
```

## Configuration
The following environment variables can be used:
- `NAMESPACE`: The Kubernetes namespace to manage pods in (default: `default`).
- `DEFAULT_REPO`: The default Docker repository for CLI images (default: `ai-cli-suite`).

## Local Development
To run locally, you need a valid `kubeconfig` and permission to manage pods. However, this manager is designed to run **in-cluster** using a ServiceAccount with appropriate RBAC.
