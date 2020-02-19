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

message BUILDING GET-STORAGE-VALUE DOCKER IMAGE
docker build -f dockerfiles/get_storage_value/Dockerfile . -t makerdao/vdb-get-storage-value:$TAG

message LOGGING INTO DOCKERHUB
echo "$DOCKER_PASSWORD" | docker login --username "$DOCKER_USER" --password-stdin

message PUSHING EXECUTE DOCKER IMAGE
docker push makerdao/vdb-execute:$TAG

message PUSHING GET-STORAGE-VALUE DOCKER IMAGE
docker push makerdao/vdb-get-storage-value:$TAG

# message DEPLOYING SERVICE
