#!/bin/sh

stack=$1

if [ -z "$stack" ]; then
  echo "$0: swarm stack name is required."
  exit 1
fi

if [ ! -f .env ]; then
  echo "$0: .env file not found!"
  exit 1
fi

while IFS='=' read -r key value; do
  [ -z "$key" ] || [ "${key#\#}" != "$key" ] && continue

  value=$(echo "$value" | sed -e 's/^"//' -e 's/"$//')

  if echo "$key" | grep -Eqi 'SECRET|PASSWORD'; then
    secret_id=$(docker secret ls -q --filter name="$key")
    if [ -n "$secret_id" ]; then
      echo "[-] removing old secret: $key"
      docker secret rm "$key" >/dev/null 2>&1
    fi

    echo "$value" | docker secret create "$key" -  >/dev/null 2>&1
    docker secret inspect "$key" --format "ID: {{.ID}}, Name: {{.Spec.Name}}, Created: {{.CreatedAt}}"
  else
    config_id=$(docker config ls -q --filter name="$key")
    if [ -n "$config_id" ]; then
      echo "[-] removing old config: $key"
      docker config rm "$key" >/dev/null 2>&1
    fi

    echo "$value" | docker config create "$key" - >/dev/null 2>&1
    docker config inspect "$key" --format "ID: {{.ID}}, Name: {{.Spec.Name}}, Created: {{.CreatedAt}}"
  fi
done < .env

echo "[*] Secrets and configs loaded successfully."
