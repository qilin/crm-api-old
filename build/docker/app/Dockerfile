FROM ubuntu:18.04

WORKDIR /app

RUN apt-get update && apt-get install -y tzdata bash ca-certificates 

COPY ./build/docker/app/files /
COPY ./api api
COPY ./assets assets
COPY ./configs configs
COPY ./web web
COPY ./bin bin
COPY ./test/testdata/plugins/so plugins


RUN chmod +x /docker/bin/* \
    && chmod +x bin

ENTRYPOINT ["/docker/bin/entrypoint.sh"]