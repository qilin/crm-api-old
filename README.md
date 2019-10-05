
## How to start on Windows
* Install **MSYS2** https://www.msys2.org/ and append `msys64\usr\bin` directory of MSYS2 to the `PATH` environment variable
* Launch **msys2.exe** application from msys64 folder
* Execute command `pacman -S make mingw64/mingw-w64-x86_64-gcc`, then you can close console and use your preffered terminal
* Install `docker`
* Install `go` and see **Development** section below

## How to start on Mac, Linux
* Install `make`, `docker`
* Install `go` and see **Development** section below

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

## Migration
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
* APP_POSTGRES_DSN - example `postgres://qilin:insecure@localhost:5567?sslmode=disable`
