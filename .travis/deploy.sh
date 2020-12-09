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
EXECUTE_NAME=vdb-mcd-execute
elif [ "$ENVIRONMENT" == "staging" ]; then
TAG=staging
EXECUTE_NAME=vdb-execute
else
   message UNKNOWN ENVIRONMENT
fi

if [ -z "$ENVIRONMENT" ]; then
    echo 'You must specify an environment (bash deploy.sh <ENVIRONMENT>).'
    echo 'Allowed values are "staging" or "prod"'
    exit 1
fi

message BUILDING EXECUTE DOCKER IMAGE
docker build -f dockerfiles/execute/Dockerfile . -t makerdao/$EXECUTE_NAME:$TAG

message BUILDING EXTRACT-DIFFS DOCKER IMAGE
docker build -f dockerfiles/extract_diffs/Dockerfile . -t makerdao/vdb-extract-diffs:$TAG

message BUILDING BACKFILL-STORAGE DOCKER IMAGE
docker build -f dockerfiles/backfill_storage/Dockerfile . -t makerdao/vdb-backfill-storage:$TAG

message BUILDING BACKFILL-EVENTS DOCKER IMAGE
docker build -f dockerfiles/backfill_events/Dockerfile . -t makerdao/vdb-backfill-events:$TAG

message LOGGING INTO DOCKERHUB
echo "$DOCKER_PASSWORD" | docker login --username "$DOCKER_USER" --password-stdin

message PUSHING EXECUTE DOCKER IMAGE
docker push makerdao/$EXECUTE_NAME:$TAG

message PUSHING EXTRACT-DIFFS DOCKER IMAGE
docker push makerdao/vdb-extract-diffs:$TAG

message PUSHING BACKFILL-STORAGE DOCKER IMAGE
docker push makerdao/vdb-backfill-storage:$TAG

message PUSHING BACKFILL-EVENTS DOCKER IMAGE
docker push makerdao/vdb-backfill-events:$TAG

# service deploy
if [ "$ENVIRONMENT" == "prod" ]; then
  message DEPLOYING EXECUTE
  aws ecs update-service --cluster vdb-cluster-$ENVIRONMENT --service vdb-mcd-execute-$ENVIRONMENT --force-new-deployment --endpoint https://ecs.$PROD_REGION.amazonaws.com --region $PROD_REGION

  message DEPLOYING EXTRACT-DIFFS
  aws ecs update-service --cluster vdb-cluster-$ENVIRONMENT --service vdb-extract-diffs-eu-$ENVIRONMENT --force-new-deployment --endpoint https://ecs.$PROD_REGION.amazonaws.com --region $PROD_REGION

  message DEPLOYING EXTRACT-DIFFS-US
  aws ecs update-service --cluster vdb-cluster-$ENVIRONMENT --service vdb-extract-diffs-us-$ENVIRONMENT --force-new-deployment --endpoint https://ecs.$PROD_REGION.amazonaws.com --region $PROD_REGION
elif [ "$ENVIRONMENT" == "staging" ]; then
  message DEPLOYING EXECUTE
  aws ecs update-service --cluster vdb-cluster-$ENVIRONMENT --service vdb-execute-$ENVIRONMENT --force-new-deployment --endpoint https://ecs.$STAGING_REGION.amazonaws.com --region $STAGING_REGION

  message DEPLOYING EXTRACT-DIFFS
  aws ecs update-service --cluster vdb-cluster-$ENVIRONMENT --service vdb-extract-diffs-$ENVIRONMENT --force-new-deployment --endpoint https://ecs.$STAGING_REGION.amazonaws.com --region $STAGING_REGION

  message DEPLOYING EXTRACT-DIFFS2
  aws ecs update-service --cluster vdb-cluster-$ENVIRONMENT --service vdb-extract-diffs2-$ENVIRONMENT --force-new-deployment --endpoint https://ecs.$STAGING_REGION.amazonaws.com --region $STAGING_REGION
else
   message UNKNOWN ENVIRONMENT
fi
