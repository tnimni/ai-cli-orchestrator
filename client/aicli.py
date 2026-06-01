#!/usr/bin/env python3
import argparse
import requests
import sys
import os
import json
import base64
import subprocess

DEFAULT_BASE_URL = os.getenv("AICLI_BASE_URL", "http://localhost:8080")

def handle_interact(args):
    # Determine the command to run based on the pod labels
    try:
        # Get pod info to find the CLI type
        pod_info = subprocess.check_output(["kubectl", "get", "pod", args.name, "-o", "json"]).decode("utf-8")
        pod_data = json.loads(pod_info)
        cli_type = pod_data.get("metadata", {}).get("labels", {}).get("cli", "")
        
        if not cli_type:
            print(f"Warning: Could not determine CLI type for pod {args.name}. Defaulting to 'bash'.")
            cmd = "bash"
        else:
            # Map type to binary
            cmd_map = {
                "claude": "claude",
                "gemini": "gemini",
                "codex": "codex",
                "bob": "bob"
            }
            cmd = cmd_map.get(cli_type, "bash")

        print(f"Connecting to {args.name} ({cmd})...")
        # Use execvp to replace current process with kubectl
        os.execvp("kubectl", ["kubectl", "exec", "-it", args.name, "--", cmd])
    except subprocess.CalledProcessError:
        print(f"Error: Could not find pod {args.name}", file=sys.stderr)
        sys.exit(1)
    except Exception as e:
        print(f"Error during interaction: {e}", file=sys.stderr)
        sys.exit(1)

def handle_create(args):
    url = f"{args.base_url}/pods"
    
    google_auth_json = ""
    if args.google_auth_file:
        with open(args.google_auth_file, "rb") as f:
            google_auth_json = base64.b64encode(f.read()).decode("utf-8")

    payload = {
        "name": args.name,
        "type": args.type,
        "repo": args.repo,
        "apiKey": args.api_key,
        "googleAuthJSON": google_auth_json,
        "useVertexAI": args.use_vertex_ai,
        "useGCA": args.use_gca,
        "mountKeys": args.mount_keys
    }
    try:
        response = requests.post(url, json=payload)
        response.raise_for_status()
        print(f"Success: {response.json()['message']}")
    except requests.exceptions.RequestException as e:
        print(f"Error creating pod: {e}", file=sys.stderr)
        if hasattr(e.response, 'text'):
            print(e.response.text, file=sys.stderr)
        sys.exit(1)

def handle_restart(args):
    url = f"{args.base_url}/pods/{args.name}/restart"
    try:
        response = requests.put(url)
        response.raise_for_status()
        print(f"Success: {response.json()['message']}")
    except requests.exceptions.RequestException as e:
        print(f"Error restarting pod: {e}", file=sys.stderr)
        sys.exit(1)

def handle_update(args):
    url = f"{args.base_url}/pods/{args.name}"
    payload = {
        "type": args.type,
        "repo": args.repo
    }
    try:
        response = requests.put(url, json=payload)
        response.raise_for_status()
        print(f"Success: {response.json()['message']}")
    except requests.exceptions.RequestException as e:
        print(f"Error updating pod: {e}", file=sys.stderr)
        sys.exit(1)

def handle_delete(args):
    url = f"{args.base_url}/pods/{args.name}"
    try:
        response = requests.delete(url)
        response.raise_for_status()
        print(f"Success: {response.json()['message']}")
    except requests.exceptions.RequestException as e:
        print(f"Error deleting pod: {e}", file=sys.stderr)
        sys.exit(1)

def main():
    parser = argparse.ArgumentParser(description="AI CLI Orchestrator Client")
    parser.add_argument("--base-url", default=DEFAULT_BASE_URL, help="Base URL of the manager service")
    subparsers = parser.add_subparsers(dest="command", required=True)

    # Create command
    create_parser = subparsers.add_parser("create", help="Create a new AI CLI pod")
    create_parser.add_argument("type", choices=["claude", "codex", "gemini", "bob"], help="Type of CLI")
    create_parser.add_argument("name", help="Name of the pod")
    create_parser.add_argument("--repo", help="Optional Docker repo override")
    create_parser.add_argument("--api-key", help="API Key for the CLI tool (e.g., GEMINI_API_KEY for gemini)")
    create_parser.add_argument("--google-auth-file", help="Path to Google Service Account JSON file")
    create_parser.add_argument("--use-vertex-ai", action="store_true", help="Use Vertex AI for Gemini")
    create_parser.add_argument("--use-gca", action="store_true", help="Use GCA for Gemini")
    create_parser.add_argument("--mount-keys", action="store_true", help="Mount local SSH keys to the pod")
    create_parser.set_defaults(func=handle_create)

    # Restart command
    restart_parser = subparsers.add_parser("restart", help="Restart an existing pod")
    restart_parser.add_argument("name", help="Name of the pod")
    restart_parser.set_defaults(func=handle_restart)

    # Update command
    update_parser = subparsers.add_parser("update", help="Update an existing pod")
    update_parser.add_argument("name", help="Name of the pod")
    update_parser.add_argument("--type", choices=["claude", "codex", "gemini", "bob"], required=True, help="New type of CLI")
    update_parser.add_argument("--repo", help="New Docker repo override")
    update_parser.set_defaults(func=handle_update)

    # Delete command
    delete_parser = subparsers.add_parser("delete", help="Delete a pod")
    delete_parser.add_argument("name", help="Name of the pod")
    delete_parser.set_defaults(func=handle_delete)

    # Interact command
    interact_parser = subparsers.add_parser("interact", help="Interact with a running pod as if local")
    interact_parser.add_argument("name", help="Name of the pod")
    interact_parser.set_defaults(func=handle_interact)

    args = parser.parse_args()
    args.func(args)

if __name__ == "__main__":
    main()
