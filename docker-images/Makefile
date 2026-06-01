# Docker repository default (can be overridden via command line: make build DOCKER_REPO=myrepo)
DOCKER_REPO ?= ai-cli-suite
TAG ?= latest

# List of CLIs
CLIS = claude codex gemini bob

.PHONY: all build push $(CLIS) clean help setup

all: build

setup:
	./setup.sh

# Build all images
build: $(CLIS)

# Generic build rule for each CLI
$(CLIS):
	@echo "Building image for $@..."
	docker build -t $(DOCKER_REPO)/$@-cli:$(TAG) -f Dockerfile.$@ .

# Push all images to the repository
push:
	@echo "Pushing images to $(DOCKER_REPO)..."
	@for cli in $(CLIS); do \
		docker push $(DOCKER_REPO)/$$cli-cli:$(TAG); \
	done

# Clean local images
clean:
	@for cli in $(CLIS); do \
		docker rmi $(DOCKER_REPO)/$$cli-cli:$(TAG) || true; \
	done

help:
	@echo "Available commands:"
	@echo "  make build          - Build all Docker images"
	@echo "  make push           - Push all images to DOCKER_REPO (default: $(DOCKER_REPO))"
	@echo "  make <cli>          - Build a specific CLI image (claude, codex, gemini, bob)"
	@echo "  make clean          - Remove local images"
	@echo ""
	@echo "Example: make build DOCKER_REPO=myregistry.com/myuser TAG=v1.0"
