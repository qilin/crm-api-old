#!/usr/bin/env sh

if [ -n "$1" ] && [ ${0:0:4} = "/bin" ]; then
  ROOT_DIR=$1/..
else
  ROOT_DIR="$( cd "$( dirname "$0" )" && pwd )/.."
fi

GO_PATH=${ROOT_DIR}/../../../..
GO_IMAGE=p1hub/qilin-docker-go
GO_IMAGE_TAG=1.12
GO_PKG=github.com/qilin/go-blueprint
GOOS="linux"
GOARCH="amd64"
DOCKER_NETWORK="blueprint-default"
DOCKER_IMAGE=p1hub/qilin-go-blueprint
PROJECT_NAME="blueprint"