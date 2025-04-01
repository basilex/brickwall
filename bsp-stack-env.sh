#!/bin/bash

ENV_FILE=".env"

SECRETS_FILE="bsp-stack-secrets.env"
CONFIGS_FILE="bsp-stack-configs.env"

echo "# Secrets file" > "$SECRETS_FILE"
echo "# Configs file" > "$CONFIGS_FILE"

while IFS= read -r line; do
    if [[ -z "$line" || "$line" =~ ^# ]]; then
        continue
    fi
    
    key="${line%%=*}"
    value="${line#*=}"
    
    if [[ "$key" =~ SECRET || "$key" =~ PASSWORD ]]; then
        echo "$line" >> "$SECRETS_FILE"
    else
        echo "$line" >> "$CONFIGS_FILE"
    fi
done < "$ENV_FILE"

echo "Files $SECRETS_FILE and $CONFIGS_FILE has been created"