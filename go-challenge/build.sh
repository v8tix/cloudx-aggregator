#!/usr/bin/env bash

IMAGE="wsclient"
TAG="v1.0.0"
IMAGE_TAG="${IMAGE}:${TAG}"

#docker build --no-cache -t "${IMAGE_TAG}" .
docker build --no-cache --progress=plain -t "${IMAGE_TAG}" .&> build.log