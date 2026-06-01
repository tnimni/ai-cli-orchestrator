# AI CLI Docker Images

This directory contains the Dockerfiles for the AI CLI tools managed by the AI CLI Orchestrator. Each image is pre-configured with essential development tools like `git`, `jq`, `yq`, and `tar`.

## Included CLIs
- **Claude CLI** (`@anthropic-ai/claude-code`)
- **Codex CLI** (`@openai/codex`)
- **Gemini CLI** (`@google/gemini-cli`)
- **Bob CLI** (`cli-bob`)

## Build & Push

While you can build these images individually, it is recommended to use the **root Makefile** from the project root:

### From Project Root (Recommended)
```bash
make build-images DOCKER_REPO=your-repo
```

### Individual Build
```bash
make build DOCKER_REPO=your-repo
```
