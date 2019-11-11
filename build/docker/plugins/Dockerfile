# Dockerfile References: https://docs.docker.com/engine/reference/builder/

## BUILDER
# Start from the latest golang base image
FROM golang:latest as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
#RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN env CGO_ENABLED=1 go build -ldflags="-w -s" -buildmode=plugin -v -o /app/gamenet.so /app/test/testdata/plugins/meta/gamenet.go
RUN env CGO_ENABLED=1 go build -ldflags="-w -s" -buildmode=plugin -v -o /app/rambler.so /app/test/testdata/plugins/parent/rambler.go
RUN pwd
RUN ls -lah

## CONTAINER
FROM ubuntu:18.04

WORKDIR /app

RUN apt-get update && apt-get install -y tzdata bash ca-certificates

RUN mkdir plugins
COPY --from=builder /app/gamenet.so plugins
COPY --from=builder /app/rambler.so plugins
COPY ./build/docker/app/files /
COPY ./api api
COPY ./assets assets
COPY ./configs configs
COPY ./bin bin
RUN pwd
RUN ls -lah
RUN ls -lah /app/plugins


RUN chmod +x /docker/bin/* \
    && chmod +x bin

ENTRYPOINT ["/docker/bin/entrypoint.sh"]