#!/bin/bash

# Detect Operating System
OS="$(uname -s)"

echo "Detected OS: $OS"

case "$OS" in
    Linux*)
        echo "Running on Linux. Ensuring Docker is available and check permissions..."
        if ! command -v docker &> /dev/null; then
            echo "Docker not found. Please install Docker for your distribution."
        else
            echo "Docker found: $(docker --version)"
        fi
        # Example of a linux specific change:
        # sudo usermod -aG docker $USER
        ;;
    Darwin*)
        echo "Running on macOS. Ensuring Docker Desktop is running..."
        if ! command -v docker &> /dev/null; then
            echo "Docker not found. Please install Docker Desktop for Mac."
        else
            echo "Docker found: $(docker --version)"
        fi
        ;;
    *)
        echo "Running on an unsupported OS: $OS"
        exit 1
        ;;
esac

echo "Setup complete. You can now use 'make build' to create the images."
