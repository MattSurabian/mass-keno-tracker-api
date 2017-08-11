#!/usr/bin/env bash

docker rm -f mass-keno-redis || exit 0
rm ../volumes/redis/data/*