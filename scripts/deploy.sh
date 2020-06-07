#!/bin/bash

NAME="carousel"
PORT="5100"
URL="root@mashdrop.com"
UPLOAD_PATH="/root/carousel/api/upload"

IMAGE="$NAME-api"
FILE="$IMAGE.tar"
docker build --no-cache -t $IMAGE .
docker save -o $FILE $IMAGE
scp $FILE $URL:~/
ssh $URL \
    "
    docker stop $NAME &&\
    docker rmi $IMAGE &&\
    docker load -i $IMAGE.tar &&\
    docker run --rm -d -p $PORT:80 -v $UPLOAD_PATH:/app/upload --name $NAME $IMAGE &&\
    rm $FILE
    "
rm $FILE