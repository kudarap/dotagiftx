#!/bin/bash

# mysql database
docker run --rm --name dota2giftables-rethinkdb \
  -p 28015:28015 -p 8080:8080 \
  -v "$PWD/.localdata/rethinkdb":/data \
  -d kudarap/rethinkdb:2.4