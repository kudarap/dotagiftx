#!/bin/bash

# requires authentication:
# podman login docker.io -u YOUR_USERNAME -p YOUR_TOKEN

image=docker.io/kudarap/dotagiftx:$1

podman build -t $image .
podman push $image