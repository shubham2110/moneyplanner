#!/bin/bash

###
## crontab -e
## */5 * * * * ./auto_update.sh >> /var/log/auto_update.log 2>&1
date
# CONFIG — change these
REPO_DIR="$(dirname "${BASH_SOURCE[0]}")"
COMPOSE_FILE="docker-compose.yml"
BRANCH="master"

# Ensure script fails on errors
set -e


cd "$REPO_DIR" || exit 1

# Fetch remote changes
git fetch origin "$BRANCH"

# Check if there are changes
LOCAL_HASH=$(git rev-parse HEAD)
REMOTE_HASH=$(git rev-parse "origin/$BRANCH")

if [ "$LOCAL_HASH" = "$REMOTE_HASH" ]; then
    # No changes
    exit 0
fi

# Pull changes
if ! git pull origin "$BRANCH"; then
    echo "Git pull failed. Skipping docker restart."
    exit 1
fi


docker-compose  -f "$COMPOSE_FILE" stop
docker-compose  -f "$COMPOSE_FILE" rm -f
# Build updated images
docker-compose -f "$COMPOSE_FILE" build

docker-compose -f "$COMPOSE_FILE" rm -f

# Restart containers with new images
docker-compose -f "$COMPOSE_FILE" up -d

exit 0
