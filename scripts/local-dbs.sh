#!/bin/bash

# rethink database
docker run --rm --name dotagiftx-rethinkdb \
  -p 28015:28015 -p 8080:8080 \
  -v "$PWD/.localdata/rethinkdb":/data \
  -d kudarap/rethinkdb:2.4


# redis cache
docker run --rm --name dotagiftx-redis \
  -p 6379:6379 \
  -e REDIS_PASSWORD=root \
  -v "$PWD/.localdata/redis":/data \
  -d bitnami/redis:6.2
