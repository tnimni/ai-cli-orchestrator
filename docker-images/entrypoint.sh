#!/bin/bash

# Decode Google Auth JSON if provided
if [ -n "$GOOGLE_AUTH_JSON" ]; then
    echo "Decoding GOOGLE_AUTH_JSON..."
    echo "$GOOGLE_AUTH_JSON" | base64 -d > /tmp/google-auth.json
    export GOOGLE_APPLICATION_CREDENTIALS=/tmp/google-auth.json
fi

if [ $# -gt 0 ]; then
    exec "$@"
else
    echo "No arguments provided. Keeping container alive for interactive use..."
    tail -f /dev/null
fi
