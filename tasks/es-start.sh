#!/usr/bin/env bash

NAME=mass-keno-es
IMAGE=docker.elastic.co/elasticsearch/elasticsearch:5.4.1
LOCAL_ES_CONF=$PWD/../volumes/elasticsearch/elasticsearch.yml
LOCAL_ES_SNAPSHOTS=$PWD/../volumes/elasticsearch/snapshots/
CONTAINER_ES_SNAPSHOTS=/usr/share/elasticsearch/snapshots/
CONTAINER_ES_CONF=/usr/share/elasticsearch/config/elasticsearch.yml

if [ "$(docker ps -aq -f status=running -f name=$NAME)" ]; then
  echo "Found running $NAME server! Run rm-es to destroy. Skipping..."
  exit 0;
fi

if [ "$(docker ps -aq -f status=exited -f name=$NAME)" ]; then
  echo "Found shut down $NAME server. Restarting..."
  docker start $NAME
  echo "Run rm-es to destroy..."
else
  echo "Starting new $NAME server..."
  docker                                         \
  run                                            \
  -d                                             \
  -p 9200:9200                                   \
  -v $LOCAL_ES_CONF:$CONTAINER_ES_CONF           \
  -v $LOCAL_ES_SNAPSHOTS:$CONTAINER_ES_SNAPSHOTS \
  -e "http.host=0.0.0.0"                         \
  -e "transport.host=127.0.0.1"                  \
  --name $NAME                                   \
  --network mass-keno                            \
  $IMAGE
fi

echo "Elastic Search is now running on localhost:9200"