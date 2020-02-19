#! /usr/bin/env bash

set -e

function message() {
    echo
    echo -----------------------------------
    echo "$@"
    echo -----------------------------------
    echo
}

ENVIRONMENT=$1
if [ "$ENVIRONMENT" == "prod" ]; then
TAG=latest
elif [ "$ENVIRONMENT" == "staging" ]; then
TAG=staging
else
   message UNKNOWN ENVIRONMENT
fi

if [ -z "$ENVIRONMENT" ]; then
    echo 'You must specifiy an envionrment (bash deploy.sh <ENVIRONMENT>).'
    echo 'Allowed values are "staging" or "prod"'
    exit 1
fi

message BUILDING EXECUTE DOCKER IMAGE
docker build -f dockerfiles/execute/Dockerfile . -t makerdao/vdb-execute:$TAG

message LOGGING INTO DOCKERHUB
echo "$DOCKER_PASSWORD" | docker login --username "$DOCKER_USER" --password-stdin

message PUSHING DOCKER IMAGE
docker push makerdao/vdb-execute:$TAG .

# message DEPLOYING SERVICE
