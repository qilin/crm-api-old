## How to start on Windows
* Install **MSYS2** https://www.msys2.org/ and append `msys64\usr\bin` directory of MSYS2 to the `PATH` environment variable
* Launch **msys2.exe** application from msys64 folder
* Execute command `pacman -S make mingw64/mingw-w64-x86_64-gcc`, then you can close console and use your preffered terminal
* Install `docker`
* Install `go` and see **Development** section below

## How to start on Mac, Linux
* Install `make`, `docker`
* Install `go` and see **Development** section below

## Quick start
* Run `make dev-docker-compose-up`
* Examine endpoints

## Endpoints
GraphQL endpoint: `http://localhost:8081/v1/graphql`  
Hasura console: `http://localhost:8081/console` (usage only as sandbox) 

## Development
* For help, run `make`
* Init local env `make up` should run only once
* Generate source files from resource `make generate`
* Build and run application `make dev-build-up` (also it's usage for rebuild && recreate containers)
* Hasura CLI for managing projects and migrations ([Install guide](https://docs.hasura.io/1.0/graphql/manual/hasura-cli/install-hasura-cli.html))

## Hasura

* `cd ./assets/hasura`
* `hasura console`
* As you use the Hasura console UI to make changes to your schema, migration files are automatically generated in the `./assets/hasura/migrations/`
* [Docs about migration & metadata](https://docs.hasura.io/1.0/graphql/manual/migrations/index.html)

## Requirements
* Docker + Docker-Compose
* Git

