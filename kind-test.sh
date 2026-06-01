#!/bin/bash
set -e

CLUSTER_NAME="ai-cli-cluster"
DOCKER_REPO="ai-cli-suite"
TAG="latest"

echo "--- Setting up Kind Cluster: $CLUSTER_NAME ---"
if kind get clusters | grep -q "^$CLUSTER_NAME$"; then
    echo "Cluster already exists."
else
    kind create cluster --name "$CLUSTER_NAME"
fi

echo "--- Building Images ---"
make build-images DOCKER_REPO=$DOCKER_REPO TAG=$TAG
make build-manager DOCKER_REPO=$DOCKER_REPO TAG=$TAG

echo "--- Loading Images into Kind ---"
kind load docker-image "$DOCKER_REPO/claude-cli:$TAG" --name "$CLUSTER_NAME"
kind load docker-image "$DOCKER_REPO/codex-cli:$TAG" --name "$CLUSTER_NAME"
kind load docker-image "$DOCKER_REPO/gemini-cli:$TAG" --name "$CLUSTER_NAME"
kind load docker-image "$DOCKER_REPO/bob-cli:$TAG" --name "$CLUSTER_NAME"
kind load docker-image "$DOCKER_REPO/ai-cli-manager:$TAG" --name "$CLUSTER_NAME"

echo "--- Deploying Chart ---"
make deploy-chart DOCKER_REPO=$DOCKER_REPO TAG=$TAG

echo "--- Waiting for Manager to be ready ---"
kubectl rollout status deployment/ai-cli-orchestrator --timeout=90s

echo "--- Setup Port-Forward ---"
echo "Run this in a separate terminal:"
echo "kubectl port-forward svc/ai-cli-orchestrator 8080:8080"

echo "--- Testing with Client CLI ---"
cd client
if [ ! -d "venv" ]; then
    python3 -m venv venv
fi
source venv/bin/activate
pip install -q requests

echo "Wait 5 seconds for port-forward (assuming you started it)..."
sleep 5

echo "Testing pod creation..."
python aicli.py create claude test-pod

echo "Verifying pod in K8s..."
kubectl get pods | grep test-pod

echo "--- Test Complete ---"
