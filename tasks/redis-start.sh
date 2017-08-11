#!/usr/bin/env bash

NAME=mass-keno-redis
IMAGE=redis
LOCAL_PORT=36379
LOCAL_REDIS_DATA=$PWD/../volumes/redis/data
CONTAINER_REDIS_DATA=/data
LOCAL_REDIS_CONF=$PWD/../volumes/redis/redis.conf
CONTAINER_REDIS_CONF=/usr/local/etc/redis/redis.conf

if [ "$(docker ps -aq -f status=running -f name=$NAME)" ]; then
  echo "Found running $NAME server! Run rm-redis to destroy. Skipping..."
  exit 0;
fi

if [ "$(docker ps -aq -f status=exited -f name=$NAME)" ]; then
  echo "Found shut down $NAME server. Restarting..."
  docker start $NAME
  echo "Run rm-redis to destroy..."
else
  echo "Starting new $NAME server..."
  docker                                       \
  run                                          \
  -d                                           \
  -p $LOCAL_PORT:6379                          \
  -v $LOCAL_REDIS_CONF:$CONTAINER_REDIS_CONF   \
  -v $LOCAL_REDIS_DATA:$CONTAINER_REDIS_DATA   \
  --name $NAME                                 \
  --network mass-keno                          \
  $IMAGE                                       \
  redis-server $CONTAINER_REDIS_CONF
fi

echo "Redis is now running on localhost:$LOCAL_PORT"