### Dev Quick Start

Requirements:

 - docker, docker-compose
 - git

To build and run project:

	docker-compose up

This automatically applies database migrations up to latest version
GraphQL endpoint: `http://localhost:8081/v1/graphql`

Reset database:

	docker-compose down -v

Hasura console ([requires hasura-cli installed](https://docs.hasura.io/1.0/graphql/manual/hasura-cli/install-hasura-cli.html#install-hasura-cli)):

	cd hasura
	hasura console

