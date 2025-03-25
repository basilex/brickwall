#!/bin/bash
#
# Docker Swarm environment loader
#
set -e
#
# Does the Swarm is working?
#
if ! docker info | grep -q "Swarm: active"; then
    echo "Docker Swarm is not initialized! Run 'docker swarm init'"
    exit 1
fi
#
# Load an environment from .env file
#
if [ ! -f .env ]; then
    echo "Environment file <.env> not found!"
    exit 1
fi

echo "Loading <.env> variables into the Docker Swarm..."
#
# Function for config creation
#
create_config() {
    key=$1
    value=$2
    config_name="config_${key}"

    echo "$value" | docker config create "$config_name" - >/dev/null 2>&1 || {
        docker config rm "$config_name" >/dev/null 2>&1
        echo "$value" | docker config create "$config_name" -
    }
    echo "> Docker Swarm <$config_name> loaded successfully"
}
#
# Function for secret creation
#
create_secret() {
    key=$1
    value=$2
    secret_name="secret_${key}"

    echo "$value" | docker secret create "$secret_name" - >/dev/null 2>&1 || {
        docker secret rm "$secret_name" >/dev/null 2>&1
        echo "$value" | docker secret create "$secret_name" -
    }
    echo "! Docker Swarm <$secret_name> loaded successfully"
}
#
# Walking through each var in the <.env> file
#
while IFS='=' read -r key value; do
    #
    # Ignore comments and empty strings
    #
    if [[ -z "$key" || "$key" == \#* ]]; then
        continue
    fi
    #
    # Determine target for writing (config or secret)
    #
    if [[ "$key" == *"_SECRET" || "$key" == *"_PASSWORD" ]]; then
        create_secret "$key" "$value"
    else
        create_config "$key" "$value"
    fi
done < .env

