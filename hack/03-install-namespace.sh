#!/bin/bash

set -euxo pipefail

PW="$1"


kubectl create namespace frank

echo "$PW"

kubectl -n frank create secret docker-registry docker-secret \
    --docker-server=registry.gitlab.com \
    --docker-username=seantywork \
    --docker-password="$1" \