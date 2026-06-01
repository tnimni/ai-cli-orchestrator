# AI CLI Orchestrator Client

A Python CLI tool to manage AI CLI pods via the Manager API. This tool is the primary interface for users to interact with the AI CLI Orchestrator Suite.

## Requirements
- Python 3.x
- `requests` library

```bash
pip install requests
```

## Configuration
You can set the base URL of the manager service using the `AICLI_BASE_URL` environment variable or the `--base-url` flag.

```bash
export AICLI_BASE_URL=http://localhost:8080
```

## Usage

### Create a Pod
```bash
python aicli.py create claude my-claude-pod
```

### Restart a Pod
```bash
python aicli.py restart my-claude-pod
```

### Update a Pod
```bash
python aicli.py update my-claude-pod --type gemini
```

### Delete a Pod
```bash
python aicli.py delete my-claude-pod
```

## Accessing the Manager inside Kubernetes
If you are running the client outside the cluster, you may need to use `kubectl port-forward`:

```bash
kubectl port-forward svc/ai-cli-orchestrator 8080:8080
```
Then run the client pointing to `localhost:8080`.
