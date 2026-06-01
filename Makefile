# Root Makefile for AI CLI Orchestrator Suite

DOCKER_REPO ?= ai-cli-suite
TAG ?= latest

# List of CLIs
CLIS = claude codex gemini bob

.PHONY: build-images build-manager push-all deploy-chart help kind-setup kind-clean $(CLIS) deploy-pod restart-pod update-pod delete-pod kind-refresh-manager kind-refresh-clis

help:
	@echo "AI CLI Orchestrator Suite Management"
	@echo "Usage:"
	@echo "  make build-images   - Build all AI CLI Docker images"
	@echo "  make build-manager  - Build the Go Manager Docker image"
	@echo "  make <cli>          - Build a specific CLI image (claude, codex, gemini, bob)"
	@echo "  make push-all       - Push all images to DOCKER_REPO"
	@echo "  make deploy-chart   - Deploy the Helm chart to Kubernetes"
	@echo "  make kind-setup             - Set up Kind cluster, build images, and load them"
	@echo "  make kind-refresh-manager   - Rebuild, reload manager to Kind, and restart"
	@echo "  make kind-refresh-clis      - Rebuild and reload all CLI images to Kind"
	@echo "  make kind-clean             - Delete the Kind cluster"
	@echo ""
	@echo "Pod Management (requires port-forward or AICLI_BASE_URL):"
	@echo "  make deploy-pod TYPE=claude NAME=my-pod API_KEY=sk-...         - Create a new AI CLI pod"
	@echo "  make deploy-pod TYPE=gemini NAME=my-pod GOOGLE_AUTH=sa.json    - Create with Service Account"
	@echo "  make deploy-pod TYPE=claude NAME=my-pod MOUNT_KEYS=true        - Create with SSH keys mounted"
	@echo "  make interact NAME=my-pod                                      - Enter pod interactively"
	@echo "  make restart-pod NAME=my-pod                                   - Restart a pod"

# Pod Management Wrappers
interact:
	python3 client/aicli.py interact $(NAME)

deploy-pod:
	python3 client/aicli.py create $(TYPE) $(NAME) --api-key "$(API_KEY)" --google-auth-file "$(GOOGLE_AUTH)" $(if $(filter true,$(USE_VERTEXAI)),--use-vertex-ai,) $(if $(filter true,$(USE_GCA)),--use-gca,) $(if $(filter true,$(MOUNT_KEYS)),--mount-keys,)

restart-pod:
	python3 client/aicli.py restart $(NAME)

delete-pod:
	python3 client/aicli.py delete $(NAME)

update-pod:
	python3 client/aicli.py update $(NAME) --type $(TYPE)

kind-setup:
	./kind-test.sh

kind-refresh-manager: build-manager
	kind load docker-image $(DOCKER_REPO)/ai-cli-manager:$(TAG) --name ai-cli-cluster
	kubectl rollout restart deployment/ai-cli-orchestrator
	kubectl rollout status deployment/ai-cli-orchestrator

kind-refresh-clis: build-images
	@for cli in $(CLIS); do \
		kind load docker-image $(DOCKER_REPO)/$$cli-cli:$(TAG) --name ai-cli-cluster; \
	done

kind-clean:
	kind delete cluster --name ai-cli-cluster

build-images:
	$(MAKE) -C docker-images build DOCKER_REPO=$(DOCKER_REPO) TAG=$(TAG)

# Proxy for individual CLI builds
$(CLIS):
	$(MAKE) -C docker-images $@ DOCKER_REPO=$(DOCKER_REPO) TAG=$(TAG)

build-manager:
	cd manager && docker build --progress=plain -t $(DOCKER_REPO)/ai-cli-manager:$(TAG) .

push-all:
	$(MAKE) -C docker-images push DOCKER_REPO=$(DOCKER_REPO) TAG=$(TAG)
	docker push $(DOCKER_REPO)/ai-cli-manager:$(TAG)

deploy-chart:
	helm upgrade --install ai-cli-orchestrator ./chart --set image.repository=$(DOCKER_REPO)/ai-cli-manager --set image.tag=$(TAG)
