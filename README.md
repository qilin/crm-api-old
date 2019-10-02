
## How to start on Windows
* Install **MSYS2** https://www.msys2.org/ and append `bin` directory of MSYS2 to the `PATH` environment variable
* Install `docker`
* All commands should be execute in **DIND** container run `make dind`
* See **Development** section below for examine commands of make

## How to start on Mac, Linux
* Install `make`, `docker`
* There are two way for development: 
    * Native (recommend) 
        * Install `go` and see **Development** section below
    * DIND container (slow file system on Mac)
        * Run `make dind` and see **Development** section below

## Development
* For help, run `make`
* Init local env `make up` should run only once
* Download dependencies `make vendor`
* Generate source files from resource `make generate`
* Build and run application `make dev-build-up` (also it's usage for rebuild && recreate containers)
* Navigate to http://localhost:8080/client in your browser and paste queries below from "Usage" section

## Usage
Create new item  
```
mutation{
  ms {
    new(name: "item1") {
      status
      id
    }
  }
}
```

Search items by query  
```
query{
  ms {
    search(query: "item1", cursor: {limit: 10, offset: 0, cursor:""}, order: ASC) {
      status
      id
      cursor {
        count
        limit
        offset
        cursor
      }
    }
  }
}
```

## Migration (for Windows run only in DIND)
* new: `./scripts/sql-migrate.sh new -env="local" {name}`
* up: `./scripts/sql-migrate.sh up -env="local"`
* down: `./scripts/sql-migrate.sh down -env="local"`
* redo: `./scripts/sql-migrate.sh redo -env="local"`
* skip: `./scripts/sql-migrate.sh skip -env="local"`
* status: `./scripts/sql-migrate.sh status -env="local"`

## Requirements
* GoLang 1.12+
* PostgreSQL 11+

## ENV
* APP_WD - work directory, default: application directory 
* APP_POSTGRES_DSN - example `postgres://blueprint:insecure@localhost:5567?sslmode=disable`
