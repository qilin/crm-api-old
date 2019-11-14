ifndef VERBOSE
.SILENT:
endif

override CURRENT_DIR = $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
override DOCKER_MOUNT_SUFFIX ?= consistent
override DOCKER_COMPOSE_ARGS ?= -f deployments/docker-compose/docker-compose.yml -f deployments/docker-compose/docker-compose-local.yml
override DOCKER_BUILD_ARGS ?= -f ${ROOT_DIR}/build/docker/app/Dockerfile
override GO_PATH = $(shell echo "$(GOPATH)" | cut -d';' -f1 | sed -e "s~^\(.\):~/\1~g" -e "s~\\\~/~g" -e "s~^/go:~~g" )


TAG ?= latest
CACHE_TAG ?= unknown_cache
GOOS ?= linux
GOARCH ?= amd64
CGO_ENABLED ?= 1
DIND_UID ?= 0
DING_GUID ?= 0

ifeq ($(GO111MODULE),auto)
override GO111MODULE = on
endif

ifeq ($(OS),Windows_NT)
override ROOT_DIR = $(shell echo $(CURRENT_DIR) | sed -e "s:^/./:\U&:g")
else
override ROOT_DIR = $(CURRENT_DIR)
endif

define go_docker
	if [ "${GO_PATH}" != "" ]; then VOLUME_PKG_MOD="-v /${GO_PATH}/pkg/mod:/${GO_PATH}/pkg/mod:${DOCKER_MOUNT_SUFFIX}"; fi ;\
	. ${ROOT_DIR}/scripts/common.sh ${ROOT_DIR}/scripts ;\
	docker run --rm \
		-v /${ROOT_DIR}:/${ROOT_DIR}:${DOCKER_MOUNT_SUFFIX} \
        $${VOLUME_PKG_MOD} \
		-w /${ROOT_DIR} \
		-e GO111MODULE=on \
		-e GOPATH=/${GO_PATH}:/go \
		$${GO_IMAGE}:$${GO_IMAGE_TAG} \
		sh -c 'GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=${CGO_ENABLED} TAG=${TAG} $(subst ",,${1}); if [ "${DIND_UID}" != "0" ]; then chown -R ${DIND_UID}:${DIND_GUID} ${ROOT_DIR}; fi'
endef

up: clean ## initialize required tools
	. ${ROOT_DIR}/scripts/common.sh ${ROOT_DIR}/scripts ;\
	if [ "${DIND}" != "1" ]; then \
		go get github.com/google/wire/cmd/wire@v0.3.0 && \
			go get github.com/99designs/gqlgen@v0.9.3 && \
			go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.19.1 && \
			go get github.com/vektah/dataloaden@v0.3.0 ;\
    fi;
.PHONY: up

down: dev-docker-compose-down clean ## reset to zero setting
.PHONY: down

build: ## build application
	if [ "${DIND}" = "1" ]; then \
		$(call go_docker,"make build") ;\
    else \
		. ${ROOT_DIR}/scripts/common.sh ${ROOT_DIR}/scripts ;\
        GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=${CGO_ENABLED} \
        go build -ldflags "-X $${GO_PKG}/cmd/version.appVersion=$(TAG)-$$(date -u +%Y%m%d%H%M)" -o "$(ROOT_DIR)/bin" main.go ;\
    fi;
.PHONY: build

clean: ## remove generated files, tidy vendor dependencies
	if [ "${DIND}" = "1" ]; then \
		$(call go_docker,"make clean") ;\
    else \
        go mod tidy ;\
    	rm -rf *.out generated/* vendor bin ;\
    fi;
.PHONY: clean

dev-build-up: build docker-image-cache dev-docker-compose-up ## create new build and recreate containers in docker-compose
.PHONY: dev-build-up

dev-build-test-plugins: ## build plugins for tests
	if [ "${DIND}" = "1" ]; then \
		$(call go_docker,"make dev-build-test-plugins") ;\
    else \
        go build -buildmode=plugin -v -o "$(ROOT_DIR)/test/testdata/plugins/so/rambler.so" "$(ROOT_DIR)/test/testdata/plugins/parent/rambler.go" ;\
        go build -buildmode=plugin -v -o "$(ROOT_DIR)/test/testdata/plugins/so/gamenet.so" "$(ROOT_DIR)/test/testdata/plugins/gamenet/gamenet.go" ;\
    fi;
.PHONY: dev-build-test-plugins

dev-docker-compose-down: ## stop and remove containers, networks, images, and volumes
	. ${ROOT_DIR}/scripts/common.sh ${ROOT_DIR}/scripts ;\
	TAG=${TAG} DOCKER_NETWORK=$${DOCKER_NETWORK} docker-compose -p $${PROJECT_NAME} ${DOCKER_COMPOSE_ARGS} down -v  ;\
	(docker network inspect $${DOCKER_NETWORK} &>/dev/null && \
	(echo "Delete docker network" && docker network rm $${DOCKER_NETWORK}) || echo "Docker network \"$${DOCKER_NETWORK}\" already deleted")
.PHONY: dev-docker-compose-down

dev-docker-compose-up: ## create and start containers
	. ${ROOT_DIR}/scripts/common.sh ${ROOT_DIR}/scripts ;\
	(docker network inspect $${DOCKER_NETWORK} >/dev/null && echo "Docker network \"$${DOCKER_NETWORK}\" already created") || \
	(echo "Create docker network \"$${DOCKER_NETWORK}\"" && docker network create $${DOCKER_NETWORK})  ;\
	TAG=${TAG} DOCKER_NETWORK=$${DOCKER_NETWORK} docker-compose -p $${PROJECT_NAME} ${DOCKER_COMPOSE_ARGS} up -d
.PHONY: dev-docker-compose-up

dev-docker-compose-recreate: dev-docker-compose-down ## pull newest version of containers and start containers
	. ${ROOT_DIR}/scripts/common.sh ${ROOT_DIR}/scripts ;\
	(docker network inspect $${DOCKER_NETWORK} >/dev/null && echo "Docker network \"$${DOCKER_NETWORK}\" already created") || \
	(echo "Create docker network \"$${DOCKER_NETWORK}\"" && docker network create $${DOCKER_NETWORK})  ;\
	TAG=${TAG} DOCKER_NETWORK=$${DOCKER_NETWORK} docker-compose -p $${PROJECT_NAME} ${DOCKER_COMPOSE_ARGS} pull --ignore-pull-failures && \
	TAG=${TAG} DOCKER_NETWORK=$${DOCKER_NETWORK} docker-compose -p $${PROJECT_NAME} ${DOCKER_COMPOSE_ARGS} up -d
.PHONY: dev-docker-compose-recreate

dev-test: test lint ## test application in dev env with race and lint
.PHONY: dev-test

dind: ## useful for windows
	if [ "${GO_PATH}" != "" ]; then VOLUME_PKG_MOD="-v /${GO_PATH}/pkg/mod:/${GO_PATH}/pkg/mod:${DOCKER_MOUNT_SUFFIX}"; fi ;\
	if [ "${DIND}" = "1" ]; then \
		echo "Already in DIND" ;\
    else \
	    . ${ROOT_DIR}/scripts/common.sh ${ROOT_DIR}/scripts ;\
	    docker rm -f dind-$${PROJECT_NAME} &>/dev/null ;\
	    docker run -it --rm --name dind-$${PROJECT_NAME} --privileged \
            -v //var/run/docker.sock://var/run/docker.sock:${DOCKER_MOUNT_SUFFIX} \
            -v /${ROOT_DIR}:/${ROOT_DIR}:${DOCKER_MOUNT_SUFFIX} \
			$${VOLUME_PKG_MOD} \
            -w /${ROOT_DIR} \
			-e GOPATH=${GO_PATH} \
            $${DIND_IMAGE}:$${DIND_IMAGE_TAG} ;\
    fi;
.PHONY: dind

docker-clean: ## delete previous docker image build
	. ${ROOT_DIR}/scripts/common.sh ${ROOT_DIR}/scripts ;\
	docker rmi $${DOCKER_IMAGE}:${CACHE_TAG} || true ;\
	docker rmi $${DOCKER_IMAGE}:${TAG} || true
.PHONY: docker-clean

docker-image-cache: ## build docker image and tagged as cache
	. ${ROOT_DIR}/scripts/common.sh ${ROOT_DIR}/scripts ;\
	docker build --cache-from $${DOCKER_IMAGE}:${CACHE_TAG} ${DOCKER_BUILD_ARGS} -t $${DOCKER_IMAGE}:${TAG} -t $${DOCKER_IMAGE}:${CACHE_TAG} ${ROOT_DIR}
.PHONY: docker-image-cache

docker-image: ## build docker image
	. ${ROOT_DIR}/scripts/common.sh ${ROOT_DIR}/scripts ;\
	docker build --cache-from $${DOCKER_IMAGE}:${CACHE_TAG} ${DOCKER_BUILD_ARGS} -t $${DOCKER_IMAGE}:${TAG} ${ROOT_DIR}
.PHONY: docker-image

docker-push: ## push docker image to registry
	. ${ROOT_DIR}/scripts/common.sh ${ROOT_DIR}/scripts ;\
	docker push $${DOCKER_IMAGE}:${TAG}
.PHONY: docker-push

generate: init gqlgen-generate go-generate ## execute all generators
.PHONY: generate

github-docker-image: docker-image ## build all docker images
	. ${ROOT_DIR}/scripts/common.sh ${ROOT_DIR}/scripts ;\
	docker build --cache-from $${DOCKER_IMAGE_HASURA}:${CACHE_TAG} -f ${ROOT_DIR}/build/docker/hasura/Dockerfile -t $${DOCKER_IMAGE_HASURA}:${TAG} ${ROOT_DIR}
.PHONY: github-docker-image

github-docker-push: docker-push ## push all docker images to registry
	. ${ROOT_DIR}/scripts/common.sh ${ROOT_DIR}/scripts ;\
	docker build --cache-from $${DOCKER_IMAGE_HASURA}:${CACHE_TAG} -f ${ROOT_DIR}/build/docker/hasura/Dockerfile -t $${DOCKER_IMAGE_HASURA}:${TAG} ${ROOT_DIR}
.PHONY: github-docker-push

github-build: github-docker-image docker-push docker-clean ## build application in CI
.PHONY: github-build

github-test: test-with-coverage ## test application in CI
.PHONY: github-test

github-demo-build: dev-build-test-plugins build docker-image docker-push ## build docker image for demo
.PHONY: github-demo-build

go-depends: ## view final versions that will be used in a build for all direct and indirect dependencies
	if [ "${DIND}" = "1" ]; then \
		$(call go_docker,"make go-depends") ;\
    else \
        cd $(ROOT_DIR) ;\
        go list -m all ;\
    fi;
.PHONY: go-depends

go-generate: init ## go generate
	if [ "${DIND}" = "1" ]; then \
		$(call go_docker,"make go-generate") ;\
    else \
        cd $(ROOT_DIR) ;\
        go generate $$(go list ./...) || exit 1 ;\
    fi;
.PHONY: go-generate

go-update-all: ## view available minor and patch upgrades for all direct and indirect
	if [ "${DIND}" = "1" ]; then \
		$(call go_docker,"make go-update-all") ;\
    else \
        cd $(ROOT_DIR) ;\
    	go list -u -m all ;\
    fi;
.PHONY: go-update-all

gqlgen-generate: ## generate graphql server
	if [ "${DIND}" = "1" ]; then \
		$(call go_docker,"make gqlgen-generate") ;\
    else \
        go run github.com/99designs/gqlgen -v -c $(ROOT_DIR)/configs/gqlgen.yml ;\
    fi;
.PHONY: gqlgen-generate

lint: ## execute linter
	if [ "${DIND}" = "1" ]; then \
		$(call go_docker,"make lint") ;\
    else \
        golangci-lint run --no-config --issues-exit-code=0 --deadline=30m \
          --disable-all --enable=deadcode  --enable=gocyclo --enable=golint --enable=varcheck \
          --enable=structcheck --enable=maligned --enable=errcheck --enable=dupl --enable=ineffassign \
          --enable=interfacer --enable=unconvert --enable=goconst --enable=gosec --enable=megacheck ./... ;\
    fi;
.PHONY: lint

test-with-coverage: ## test application with race and total coverage
	if [ "${DIND}" = "1" ]; then \
		$(call go_docker,"make test-with-coverage") ;\
    else \
		export WD=$(ROOT_DIR) ;\
        CGO_ENABLED=1 \
        go test -v -race -covermode atomic -coverprofile coverage.out ${TEST_ARGS} ./... || exit 1 ;\
        go tool cover -func=coverage.out ;\
    fi;
.PHONY: test-with-coverage

test: ## test application with race
	if [ "${DIND}" = "1" ]; then \
		$(call go_docker,"make test") ;\
    else \
		export WD=$(ROOT_DIR) ;\
        CGO_ENABLED=1 \
        go test -race -v ${TEST_ARGS} ./... ;\
    fi;
.PHONY: test

vendor: ## update vendor dependencies
	if [ "${DIND}" = "1" ]; then \
		$(call go_docker,"make vendor") ;\
    else \
        rm -rf $(ROOT_DIR)/vendor ;\
    	go mod vendor ;\
    fi;
.PHONY: vendor

go-download-deps: ## download dependencies
	if [ "${DIND}" = "1" ]; then \
		$(call go_docker,"make go-download-deps") ;\
    else \
    	go get -d ./... ;\
    fi;
.PHONY: go-download-deps

init:
	mkdir -p $(ROOT_DIR)/internal/generated
.PHONY: init

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.PHONY: help

fixtures:
	find ./assets/fixtures/ -type f -name '*.sql' -exec psql postgres://postgres:postgres@localhost:5567/qilin-hasura -f {} +
.PHONY: fixtures

fixtures-migrate:
	. ${ROOT_DIR}/scripts/common.sh ${ROOT_DIR}/scripts ;\
	docker run --rm --network container:qilin-postgres -ti p1hub/qilin-crm-api:${TAG} migrate up --dsn 'postgres://postgres:postgres@postgres:5432/qilin-hasura?sslmode=disable' --path /assets/fixtures	
.PHONY: fixtures-migrate

.DEFAULT_GOAL := help