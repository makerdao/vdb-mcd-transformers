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
    echo 'You must specify an environment (bash deploy.sh <ENVIRONMENT>).'
    echo 'Allowed values are "staging" or "prod"'
    exit 1
fi

message BUILDING EXECUTE DOCKER IMAGE
docker build -f dockerfiles/execute/Dockerfile . -t makerdao/vdb-execute:$TAG

message BUILDING BACKFILL-STORAGE DOCKER IMAGE
docker build -f dockerfiles/backfill_storage/Dockerfile . -t makerdao/vdb-backfill-storage:$TAG

message BUILDING BACKFILL-EVENTS DOCKER IMAGE
docker build -f dockerfiles/backfill_events/Dockerfile . -t makerdao/vdb-backfill-events:$TAG

message BUILDING EXTRACT-DIFFS DOCKER IMAGE
docker build -f dockerfiles/extract_diffs/Dockerfile . -t makerdao/vdb-extract-diffs:$TAG

message LOGGING INTO DOCKERHUB
echo "$DOCKER_PASSWORD" | docker login --username "$DOCKER_USER" --password-stdin

message PUSHING EXECUTE DOCKER IMAGE
docker push makerdao/vdb-execute:$TAG

message PUSHING BACKFILL-STORAGE DOCKER IMAGE
docker push makerdao/vdb-backfill-storage:$TAG

message PUSHING BACKFILL-EVENTS DOCKER IMAGE
docker push makerdao/vdb-backfill-events:$TAG

message PUSHING EXTRACT-DIFFS DOCKER IMAGE
docker push makerdao/vdb-extract-diffs:$TAG

# service deploy
if [ "$ENVIRONMENT" == "prod" ]; then
  message DEPLOYING EXECUTE
  aws ecs update-service --cluster vdb-cluster-$ENVIRONMENT --service vdb-execute-$ENVIRONMENT --force-new-deployment --endpoint https://ecs.$PROD_REGION.amazonaws.com --region $PROD_REGION
elif [ "$ENVIRONMENT" == "staging" ]; then
  message DEPLOYING EXECUTE
  aws ecs update-service --cluster vdb-cluster-$ENVIRONMENT --service vdb-execute-$ENVIRONMENT --force-new-deployment --endpoint https://ecs.$STAGING_REGION.amazonaws.com --region $STAGING_REGION

  message DEPLOYING BACKFILL-EVENTS
  aws ecs run-task --cluster vdb-cluster-$ENVIRONMENT \
   --launch-type FARGATE \
   --task-definition vdb-backfill-events-$ENVIRONMENT \
   --network-configuration "$STAGING_NETWORK_CONFIG" \
   --region $STAGING_REGION

   message DEPLOYING BACKFILL-STORAGE
   aws ecs run-task --cluster vdb-cluster-$ENVIRONMENT \
    --launch-type FARGATE \
    --task-definition vdb-backfill-storage-$ENVIRONMENT \
    --network-configuration "$STAGING_NETWORK_CONFIG" \
    --region $STAGING_REGION
else
   message UNKNOWN ENVIRONMENT
fi
