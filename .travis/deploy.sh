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
if [ -z "$ENVIRONMENT" ]; then
    echo 'You must specify an environment (bash deploy.sh <ENVIRONMENT>).'
    echo 'Allowed values are "staging" or "prod"'
    exit 1
fi

#--------------------------
# INIT
#--------------------------
if [ "$ENVIRONMENT" == "prod" ]; then
    TAG=latest
    REGION=$PROD_REGION
    ECS_NETWORK_CONFIG=$PROD_NETWORK_CONFIG
elif [ "$ENVIRONMENT" == "staging" ]; then
    TAG=staging
    REGION=$STAGING_REGION
    ECS_NETWORK_CONFIG=$STAGING_NETWORK_CONFIG
else
    message UNKNOWN ENVIRONMENT
fi

#--------------------------
# BUILD IMAGES
#--------------------------
message BUILD DOCKER IMAGES

message BUILDING EXTRACT-DIFFS DOCKER IMAGE
docker build -f dockerfiles/extract_diffs/Dockerfile . -t makerdao/vdb-extract-diffs:$TAG

message BUILDING BACKFILL-STORAGE DOCKER IMAGE
docker build -f dockerfiles/backfill_storage/Dockerfile . -t makerdao/vdb-backfill-storage:$TAG

message BUILDING BACKFILL-EVENTS DOCKER IMAGE
docker build -f dockerfiles/backfill_events/Dockerfile . -t makerdao/vdb-backfill-events:$TAG

if [ "$ENVIRONMENT" == "prod" ]; then

    message BUILDING EXECUTE DOCKER IMAGE
    docker build -f dockerfiles/execute/Dockerfile . -t makerdao/vdb-mcd-execute:$TAG

elif [ "$ENVIRONMENT" == "staging" ]; then

    message BUILDING EXECUTE DOCKER IMAGE
    docker build -f dockerfiles/execute/Dockerfile . -t makerdao/vdb-execute:$TAG

else
    message UNKNOWN ENVIRONMENT
fi

#--------------------------
# DOCKERHUB
#--------------------------
message LOGGING INTO DOCKERHUB
echo "$DOCKER_PASSWORD" | docker login --username "$DOCKER_USER" --password-stdin

#--------------------------
# PUSH IMAGES
#--------------------------
message PUSH DOCKER IMAGES TO DOCKERHUB

message PUSHING BACKFILL-STORAGE DOCKER IMAGE
docker push makerdao/vdb-backfill-storage:$TAG

message PUSHING BACKFILL-EVENTS DOCKER IMAGE
docker push makerdao/vdb-backfill-events:$TAG

message PUSHING EXTRACT-DIFFS DOCKER IMAGE
docker push makerdao/vdb-extract-diffs:$TAG

if [ "$ENVIRONMENT" == "prod" ]; then

    message PUSHING EXECUTE DOCKER IMAGE
    docker push makerdao/vdb-mcd-execute:$TAG

elif [ "$ENVIRONMENT" == "staging" ]; then

    message PUSHING EXECUTE DOCKER IMAGE
    docker push makerdao/vdb-execute:$TAG

else
    message UNKNOWN ENVIRONMENT
    exit 1 # don't continue
fi

#--------------------------
# DEPLOY
#--------------------------
message DEPLOY TO $ENVIRONMENT

message DEPLOYING BACKFILL-EVENTS
aws ecs run-task \
  --cluster vdb-cluster-$ENVIRONMENT \
  --launch-type FARGATE \
  --task-definition vdb-backfill-events-$ENVIRONMENT \
  --network-configuration $ECS_NETWORK_CONFIG \
  --region $REGION

message DEPLOYING BACKFILL-STORAGE
aws ecs run-task \
  --cluster vdb-cluster-$ENVIRONMENT \
  --launch-type FARGATE \
  --task-definition vdb-backfill-storage-$ENVIRONMENT \
  --network-configuration $ECS_NETWORK_CONFIG \
  --region $REGION

if [ "$ENVIRONMENT" == "prod" ]; then
    EXECUTE_NAME = vdb-mcd-execute
    EXTRACT_DIFFS_NAME = vdb-extract-diffs-eu
    EXTRACT_DIFFS_US_NAME = vdb-extract-diffs-us
elif [ "$ENVIRONMENT" == "staging" ]; then
    EXECUTE_NAME=vdb-execute
    EXTRACT_DIFFS_NAME = vdb-extract-diffs
    EXTRACT_DIFFS_US_NAME = vdb-extract-diffs2
fi

message DEPLOYING EXECUTE
aws ecs update-service \
  --cluster vdb-cluster-$ENVIRONMENT \
  --service $EXECUTE_NAME-$ENVIRONMENT \
  --force-new-deployment \
  --endpoint https://ecs.$REGION.amazonaws.com \
  --region $REGION

message DEPLOYING EXTRACT-DIFFS
aws ecs update-service \
  --cluster vdb-cluster-$ENVIRONMENT \
  --service $EXTRACT_DIFFS_NAME-$ENVIRONMENT \
  --force-new-deployment \
  --endpoint https://ecs.$REGION.amazonaws.com \
  --region $REGION

message DEPLOYING EXTRACT-DIFFS-US
aws ecs update-service \
  --cluster vdb-cluster-$ENVIRONMENT \
  --service $EXTRACT_DIFFS_US_NAME-$ENVIRONMENT \
  --force-new-deployment \
  --endpoint https://ecs.$REGION.amazonaws.com \
  --region $REGION


# announce deploy
.travis/announce.sh $ENVIRONMENT vdb-mcd-transformers
