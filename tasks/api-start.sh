#!/usr/bin/env bash
NAME=mass-keno-api
VERSION=$(git describe --tags --always --dirty)

if [ "$(docker ps -aq -f status=running -f name=$NAME)" ]; then
  echo "Found running $NAME server! Run api-stop to destroy. Skipping..."
  exit 0;
fi

if [ "$(docker ps -aq -f status=exited -f name=$NAME)" ]; then
  echo "Found shut down $NAME server. Destroying and recreating..."
  docker rm -f $NAME
fi

docker run \
-d \
-e REDIS_CACHE_HOST=mass-keno-redis:6379 \
-e MASS_KENO_HTTP_ADDR=0.0.0.0:8090 \
-p 8090:8090 \
--network mass-keno \
--name $NAME \
mattsurabian/mass-keno-tracker-api-amd64:$VERSION

echo "Mass Keno Tracker API now running on localhost:8090/v1/health"