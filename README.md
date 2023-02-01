# Introduction
            .,*/(#####(/*,.                               .,*((###(/*.
        .*(%%%%%%%%%%%%%%#/.                           .*#%%%%####%%%%#/.
      ./#%%%%#(/,,...,,***.           .......          *#%%%#*.   ,(%%%#/.
     .(#%%%#/.                    .*(#%%%%%%%##/,.     ,(%%%#*    ,(%%%#*.
    .*#%%%#/.    ..........     .*#%%%%#(/((#%%%%(,     ,/#%%%#(/#%%%#(,
    ./#%%%(*    ,#%%%%%%%%(*   .*#%%%#*     .*#%%%#,      *(%%%%%%%#(,.
    ./#%%%#*    ,(((##%%%%(*   ,/%%%%/.      .(%%%#/   .*#%%%#(*/(#%%%#/,
     ,#%%%#(.        ,#%%%(*   ,/%%%%/.      .(%%%#/  ,/%%%#/.    .*#%%%(,
      *#%%%%(*.      ,#%%%(*   .*#%%%#*     ./#%%%#,  ,(%%%#*      .(%%%#*
       ,(#%%%%%##(((##%%%%(*    .*#%%%%#(((##%%%%(,   .*#%%%##(///(#%%%#/.
         .*/###%%%%%%%###(/,      .,/##%%%%%##(/,.      .*(##%%%%%%##(*,
              .........                ......                .......
A starter kit for Go API development. Inspired by [How I write HTTP services after eight years](https://pace.dev/blog/2018/05/09/how-I-write-http-services-after-eight-years.html).

However, I wanted to use [chi router](https://github.com/go-chi/chi) which is more common in the community, [sqlx](https://github.com/jmoiron/sqlx) for database operations and design towards layered architecture (handler -> business logic -> repository).

It is still in early stages, and I do not consider it is completed until all integration tests are completed.

In short, this kit is a Go + Postgres + Chi Router + sqlx + ent + unit testing starter kit for API development.

# Motivation

On the topic of API development, there are two opposing camps between using framework (like [echo](https://github.com/labstack/echo), [gin](https://github.com/gin-gonic/gin), [buffalo](http://gobuffalo.io/)) and starting small and only adding features you need through various libraries. 

However, the second option isn't that straightforward. you will want to structure your project in such a way that there are clear separation of functionalities for your controller, business logic and database operations. Dependencies need to be injected from outside to inside. Being modular, swapping a library like a router or database library to a different one becomes much easier.

# Features

This kit is composed of standard Go library together with some well-known libraries to manage things like router, database query and migration support.

  - [x] Framework-less and net/http compatible handler
  - [x] Router/Mux with [Chi Router](https://github.com/go-chi/chi)
  - [x] Database Operations with [sqlx](https://github.com/jmoiron/sqlx)
  - [x] Database Operations with [ent](https://entgo.io/docs/getting-started)
  - [x] Database migration with [goose](https://github.com/pressly/goose)
  - [x] Input [validation](https://github.com/go-playground/validator) that returns multiple error strings
  - [x] Read all configurations using a single `.env` file or environment variable
  - [x] Clear directory structure, so you know where to find the middleware, domain, server struct, handle, business logic, store, configuration files, migrations etc. 
  - [x] (optional) Request log that logs each user uniquely based on host address
  - [x] CORS
  - [x] Scans and auto-generate [Swagger](https://github.com/swaggo/swag) docs using a declarative comments format 
  - [x] Custom model JSON output
  - [x] Filters (input DTO), Resource (output DTO) for pagination parsing and custom response respectively.
  - [x] Cache layer
  - [x] Uses [Task](https://taskfile.dev) to simplify various tasks like mocking, linting, test coverage, hot reload etc
  - [x] Unit testing of repository, use case, and handler using mocks and [dockertest](https://github.com/ory/dockertest)
  - [ ] End-to-end test using ephemeral docker containers

# Quick Start

It is advisable to use the latest supported [Go version](#appendix) (>= v1.18). Optionally `docker` and `docker-compose` for easier start up.

Get it

```shell
git clone https://github.com/gmhafiz/go8
cd go8
```

The application depends on a database. Ideally applications load configurations from [environment variable](https://12factor.net/config) or from a vault. 

1. Set values by exporting them into environment variables

```shell
export DB_DRIVER=postgres
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=user
export DB_PASS=password
export DB_NAME=go8_db
```

It is also possible to set them in `env` file. Just make sure this file is ignored in `.gitignore` because it should never be checked into source control.

2. Fill in your database credentials in `.env` by making a copy of `env.example` first.
```shell
cp env.example .env
```

Have a database ready either by installing them yourself or the following command. The `docker-compose.yml` will use database credentials set in either `.env` file or environment variables which is initialized in the previous step. Optionally, you may want redis as well.

```shell
docker-compose up -d postgres
# or
docker-compose up -d postgres redis
```

Once the database is up you may run the migration with,

```shell
go run cmd/migrate/main.go
```

Run the API with the following command. For the first time run, dependencies will be downloaded first.

```shell
go run cmd/go8/main.go
```

You will see the address of the API is running at.

```shell
2021/10/31 10:49:11 Starting API version: v0.1.0
2021/10/31 10:49:11 Connecting to database...
2021/10/31 10:49:11 Database connected
        .,*/(#####(/*,.                               .,*((###(/*.
    .*(%%%%%%%%%%%%%%#/.                           .*#%%%%####%%%%#/.
  ./#%%%%#(/,,...,,***.           .......          *#%%%#*.   ,(%%%#/.
 .(#%%%#/.                    .*(#%%%%%%%##/,.     ,(%%%#*    ,(%%%#*.
.*#%%%#/.    ..........     .*#%%%%#(/((#%%%%(,     ,/#%%%#(/#%%%#(,
./#%%%(*    ,#%%%%%%%%(*   .*#%%%#*     .*#%%%#,      *(%%%%%%%#(,.
./#%%%#*    ,(((##%%%%(*   ,/%%%%/.      .(%%%#/   .*#%%%#(*/(#%%%#/,
 ,#%%%#(.        ,#%%%(*   ,/%%%%/.      .(%%%#/  ,/%%%#/.    .*#%%%(,
  *#%%%%(*.      ,#%%%(*   .*#%%%#*     ./#%%%#,  ,(%%%#*      .(%%%#*
   ,(#%%%%%##(((##%%%%(*    .*#%%%%#(((##%%%%(,   .*#%%%##(///(#%%%#/.
     .*/###%%%%%%%###(/,      .,/##%%%%%##(/,.      .*(##%%%%%%##(*,
          .........                ......                .......
2021/10/31 10:49:11 Serving at 0.0.0.0:3080
```

To use, open a new terminal and follow examples in the `examples/` folder

```shell
curl -v --location --request POST 'http://localhost:3080/api/v1/book' \
 --header 'Content-Type: application/json' \
 --data-raw '{
    "title": "Test title",
    "image_url": "https://example.com",
    "published_date": "2020-07-31T15:04:05.123499999Z",
    "description": 
    "test description"
  }' \
 | jq

curl --location --request GET 'http://localhost:3080/api/v1/book' | jq
```

To see all available routes, run

```shell
go run cmd/route/main.go
```

![go run cmd/routes/main.go](assets/routes.png)


To run all tests,

```shell
go test ./...
```


# Table of Contents

- [Introduction](#introduction)
- [Motivation](#motivation)
- [Features](#features)
- [Quick Start](#quick-start)
- [Tooling](#tooling)
   * [Tools](#tools)
      + [Install](#install)
   * [Tasks](#tasks)
      + [List Routes](#list-routes)
      + [Format Code](#format-code)
      + [Sync Dependencies](#sync-dependencies)
      + [Compile Check](#compile-check)
      + [Unit tests](#unit-tests)
      + [golangci Linter](#golangci-linter)
      + [Security Checks](#security-checks)
      + [Check](#check)
      + [Hot reload](#hot-reload)
      + [Generate Swagger Documentation](#generate-swagger-documentation)
      + [Go generate](#go-generate)
      + [Test Coverage](#test-coverage)
      + [Build](#build)
      + [Clean](#clean)
- [Migration](#migration)
   * [Using Task](#using-task)
      + [Create Migration](#create-migration)
      + [Migrate up](#migrate-up)
      + [Rollback](#rollback)
   * [Without Task](#without-task)
      + [Create Migration](#create-migration-1)
      + [Migrate Up](#migrate-up)
      + [Rollback](#rollback-1)
- [Run](#run)
   * [Local](#local)
   * [Docker](#docker)
      + [docker-compose](#docker-compose)
- [Build](#build-1)
   * [With Task](#with-task)
   * [Without Task](#without-task-1)
- [Swagger docs](#swagger-docs)
- [Structure](#structure)
   * [Starting Point](#starting-point)
   * [Configurations](#configurations)
      - [.env files](#env-files)
   * [Database](#database)
   * [Router](#router)
   * [Domain](#domain)
      + [Repository](#repository)
      + [Use Case](#use-case)
      + [Handler](#handler)
      + [Initialize Domain](#initialize-domain)
   * [Middleware](#middleware)
   * [Dependency Injection](#dependency-injection)
   * [Libraries](#libraries)
- [Cache](#cache)
   * [LRU](#lru)
   * [Redis](#redis)
- [Utility](#utility)
- [Testing](#testing)
   * [Unit Testing](#unit-testing)
      - [Handler](#handler-1)
      - [Use Case](#use-case-1)
      - [Repository](#repository-1)
   * [End to End Test](#end-to-end-test)
- [TODO](#todo)
- [Acknowledgements](#acknowledgements)
- [Appendix](#appendix)
   * [Dev Environment Installation](#dev-environment-installation)


# Tooling

The above quick start is sufficient to start the API. However, we can take advantage of a tool to make task management easier. While you may run migration with `go run cmd/migrate/main.go`,  it is a lot easier to remember to type `task migrate` instead. Think of it as a simplified `Makefile`.

You may also choose to run sql scripts directly from `database/migrations` folder instead.

This project uses [Task](https://github.com/go-task/task) to handle various tasks such as migration, generation of swagger docs, build and run the app. It is essentially a [sh interpreter](https://github.com/mvdan/sh).

Install task runner binary bash script:

    sudo ./scripts/install-task.sh

This installs `task` to `/usr/local/bin/task` so `sudo` is needed.

`Task` tasks are defined inside `Taskfile.yml` file. A list of tasks available can be viewed with:

    task -l   # or
    task list

## Tools

Various tooling can be installed automatically by running which includes

 * [golang-ci](https://golangci-lint.run)
    * An opinionated code linter from https://golangci-lint.run/
 * [swag](https://github.com/swaggo/swag)
    * Generates swagger documentation
 * [goose](https://github.com/pressly/goose)
    * Migration tool
 * [ent](https://entgo.io/docs/getting-started)
    * Database ORM tool
 * [mirip](https://github.com/gmhafiz/mirip)
    * Generate mocks from interface 
 * [air](https://github.com/cosmtrek/air)
    * Hot reload app 

### Install

Install the tools above with:

    task install:tools


## Tasks

Various tooling are included within the `Task` runner. Configurations are done inside `Taskfile.yml` file.

### List Routes

List all registered routes, typically done by `register.go` files by

    go run cmd/route/route.go

or

    task routes

### Go generate

    task generate

Runs `go generate ./...`. It looks for `//go:generate` tags found in .go files. Useful for recreating mock file for unit tests.

### Generate Swagger Documentation

    task swagger

Reads annotations from controller and model file to create a swagger documentation file. Can be accessed from [http://localhost:3080/swagger/](http://localhost:3080/swagger/)

### Format Code

    task fmt

Runs `go fmt ./...` to lint Go code

`go fmt` is part of official Go toolchain that formats your code into an opinionated format.

### Compile Check

    task vet

Quickly catches compile error.

### golangci Linter

    task lint

Runs [https://golangci-lint.run](https://golangci-lint.run/) linter. Includes [gosec](https://github.com/securego/gosec) via `golangci.yaml`.

### Unit tests

    task test

Runs unit tests.

### Check

    task check

Runs all the above tasks (Format Code until Security Checks)

### Sync Dependencies

    task tidy

Runs `go mod tidy` to sync dependencies.

### Hot reload

    task dev

Runs `air` which watches for file changes and rebuilds binary. Configure in `.air.toml` file.

### Test Coverage

    task coverage

Runs unit test coverage with `go test -cover ./...`

### Build

    task build

Create a statically linked executable for linux.

### Clean

    task clean

Clears all files inside `bin` directory.

# Migration

Migration is a good step towards having a versioned database and makes publishing to a production server a safe process.

All migration files are stored in `database/migrations` folder.

## Using Task

### Create Migration

Using `Task`, creating a migration file is done by the following command. Name the file after `NAME=`.

    task migrate:create NAME=create_a_tablename

Write your schema in pure sql in the 'up' section and any reversal in the 'down' section of the file.
 
### Migrate up

After you are satisfied with your `.sql` files, run the following command to migrate your database.

    task migrate

To migrate one step

    task migrate:step
      
### Rollback
    
To roll back migration by one step

    task migrate:rollback

Further `goose` commands are available in its [page](https://github.com/pressly/goose)


## Without Task

### Create Migration

Once `goose` tool is [installed](https://github.com/pressly/goose), create a migration with

    migrate create -ext sql -dir database/migrations -format unix "{{.NAME}}"

### Migrate Up

You will need to create a data source name string beforehand. e.g.:

    postgres://postgres_user:$password@$localhost:5432/db?sslmode=false

Note: You can save the above string into an environment variable for reuse e.g.

    export DSN=postgres://postgres_user:$password@$localhost:5432/db?sslmode=false

Then migrate with the following command, specifying the path to migration files, data source name and action.

    migrate -path database/migrations -database $DSN up

To migrate 1 step,

    migrate -path database/migrations -database $DSN up-by-one

### Rollback

Rollback migration by using `down` action and the number of steps

    migrate -path database/migrations -database $DSN down

# Run

## Local

Conventionally, all apps are placed inside the `cmd` folder.

If you have `Task` installed, the server can be run with:

    task run

or without `Task`, just like in quick start section:

    go run cmd/go8/main.go

## Docker

You can build a docker image with the app with its config files. Docker needs to be installed beforehand.

     task docker:build

This task also makes a copy of `.env`. Since Docker doesn't copy hidden file, we make a copy of it on our `src` stage before transferring it to our final `scratch` stage. It also inserts formats git tag and git hash as the API version which runs at compile time.

Note that this is a multistage Dockerfile. Since we statically compile this API, we can use a minimal images like a `distroless` that includes both timezones and CA certificates.

Run the following command to build a container from this image. `--net=host` tells the container to use local's network so that it can access host database.

    docker-compose up -d postgres # If you haven't run this from quick start 
    task docker:run

### docker-compose

If you prefer to use docker-compose instead, both server and the database can be run with:

    task docker-compose:start

# Build

## With Task

If you have task installed, simply run

    task build

It does task check prior to build and puts both the binary and `.env` files into `./bin` folder

## Without Task

    go mod download
    CGO_ENABLED=0 GOOS=linux
    go build -ldflags="-X main.Version=$(git describe --abbrev=0 --tags)-$(git rev-list -1 HEAD) -w -s" -o ./server ./cmd/go8/main.go;


# Swagger docs

Swagger UI allows you to play with the API from a browser

![swagger UI](assets/swagger.png)
     
Edit `cmd/go8/go8.go` `main()` function host and BasePath  

    // @host localhost:3080
    // @BasePath /api/v1

   
Generate with

    task swagger # runs: swag init 
    
Access at

    http://localhost:3080

The command `swag init` scans the whole directory and looks for [swagger's declarative comments](https://github.com/swaggo/swag#declarative-comments-format) format.

Custom theme is obtained from [https://github.com/ostranme/swagger-ui-themes](https://github.com/ostranme/swagger-ui-themes)

# Structure

This project follows a layered architecture mainly consisting of three layers:

 1. Handler
 2. Use Case
 3. Repository

![layered-architecture](assets/layered-architecture.png)

The handler is responsible to receiving requests, validating them hand over to business logic.
Values returned from use case layer is then formatted and to be returned to the client.

Business logic (use case) is the meat of operations, and it calls a repository if necessary.

Database calls lives in this repository layer where data is retrieved from a store.

All of these layers are encapsulated in a domain, and an API can contain many domain.

Each layer communicates through an interface which means the layer depends on
abstraction instead of concrete implementation. This achieves loose-coupling and 
makes unit testing easier.

## Starting Point

Starting point of project is at `cmd/go8/main.go`

![main](assets/main.png)


The `Server` struct in `internal/server/server.go` is where all important dependencies are 
registered and to give a quick glance on what your server needs.

![server](assets/server.png)

`s.Init()` in `internal/server/server.go` simply initializes server configuration, database, input validator, router, global middleware, domains, and swagger. Any new dependency added to the `Server` struct can be initialized here too.

![init](assets/init.png)


## Configurations
![configs](assets/configs.png)

All environment variables are read into specific `Configs` struct initialized in `configs/configs.go`. Each of the embedded struct are defined in its own file of the same package where its fields are read from either environment variable or `.env` file.

This approach allows code completion when accessing your configurations.

![config code completion](assets/config-code-completion.png)


#### .env files

The `.env` file defines settings for various parts of the API including the database credentials. If you choose to export the variables into environment variables for example:

    export DB_DRIVER=postgres
    export DB_HOST=localhost
    export DB_PORT=5432
    etc


To add a new type of configuration, for example for Elasticsearch
 
1. Create a new go file in `./configs`

```shell
touch configs/elasticsearch.go
```
    
2. Create a new struct for your type

```go
type Elasticsearch struct {
  Address  string
  User     string
  Password string
}
```
    
3. Add a constructor for it

```go
func ElasticSearch() Elasticsearch {
   var elasticsearch Elasticsearch
   envconfig.MustProcess("ELASTICSEARCH", &elasticsearch)

   return elasticsearch
}
``` 

A namespace is defined 

4. Add to `.env` of the new environment variables

```shell
ELASTICSEARCH_ADDRESS=http://localhost:9200
ELASTICSEARCH_USER=user
ELASTICSEARCH_PASS=password
```

Limiting the number of connection pool avoids ['time-slicing' of the CPU](https://github.com/brettwooldridge/HikariCP/wiki/About-Pool-Sizing). Use the following formula to determine a suitable number
 
    number of connections = ((core_count * 2) + effective_spindle_count)    

## Database

Migrations files are stored in `database/migrations` folder. [goose](https://github.com/pressly/goose) library is used to perform migration using `task` commands.

## Router

Router multiplexer or mux is created for use by `Domain`. While [chi](https://github.com/go-chi/chi) library is being used here, you can swap out the router to an alternative one when assigning `s.router` field. However, you will need to adjust how you register your handlers in each domain.

## Domain

Let us look at how this project attempts at layered architecture. A domain consists of: 

  1. Handler (Controllers)
  2. Use case (Business Logic)
  3. Repository (Database)

Let us start by looking at how `repository` is implemented.

### Repository

Starting with `Database`. This is where all database operations are handled. Inside the `internal/domain/health` folder:

![book-domain](assets/domain-health.png)

Interfaces for both use case and repository are on its own file under the `health` package while its implementation are in `usecase` and `repository` package respectively.

The `health` repository has a single method

`internal/domain/health/repository.go`

```go
 type Repository interface {
     Readiness() error
 }
````    

And it is implemented in a package called `postgres` in `internal/domain/health/repository/postgres/postgres.go`

```go
func (r *repository) Readiness() error {
  return r.db.Ping()
}
```

### Use Case

This is where all business logic lives. By having repository layer underneath in a separate layer, those functions are reusable in other use case layers.

### Handler

This layer is responsible in handling request from outside world and into the `use case` layer. It does the following:

 1. Parse request into private 'request' struct
 2. Sanitize and validates said struct
 3. Pass into `use case` layer
 4. Process results from coming from `use case` layer and decide how the payload is going to be formatted to the outside world.
  
Route API are defined in `RegisterHTTPEndPoints` in their respective `register.go` file. 


### Initialize Domain

Finally, a domain is initialized by wiring up all dependencies in server/initDomains.go. Here, any dependencies can be injected such as a custom logger.

```go
func (s *Server) initBook() {
   newBookRepo := bookRepo.New(s.GetDB())
   newBookUseCase := bookUseCase.New(newBookRepo)
    s.Domain.Book = bookHandler.RegisterHTTPEndPoints(s.router, newBookUseCase)
}
```

## Middleware

A middleware is just a handler that returns a handler as can be seen in the `internal/middleware/cors.go`

```go
func Cors(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    
        // do something before going into Handler
        
        next.ServerHTTP(w, r)
        
        // do something after handler has been served
    }
}
```

Then you may choose to have this middleware to affect all routes by registering it in`setGlobalMiddleware()` or only a specific domain at `RegisterHTTPEndPoints()` function in its `register.go` file. 


### Middleware External Dependency

Sometimes you need to add an external dependency to the middleware which is often the case for 
authorization be that a config or a database. To do that, we wrap our middleware with a
`func(http.Handler) http.Handler`. Any dependencies can now be passed in into `Auth()`.

```go
func Auth(cfg configs.Configs) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            claims, err := getClaims(r, cfg.Jwt.SecretKey)
            if err != nil {
                w.WriteHeader(http.StatusUnauthorized)
                return
            }
    
            next.ServeHTTP(w, r)
        })
    }
}
```

## Dependency Injection

Dependency injection in Go is simple. We can simply pass in whatever we need
into the function or method signature. There is no need to make a dependency injection
container like other object-oriented programming language.

How does dependency injection happens? It starts with `InitDomains()` method. We
initialize the dependency we want early in the start-up of the program and then 
pass it down the layers.

```go
newBookRepo := bookRepo.New(s.DB())
newBookUseCase := bookUseCase.New(newBookRepo)
s.Domain.Book = bookHandler.RegisterHTTPEndPoints(
	s.router,
	s.validator,
	newBookUseCase,
)
```

The repository gets access to a pointer to `sql.DB` from the `s.DB()` 
initialisation to perform database operations. This layer also knows nothing of
layers above it. `NewBookUseCase` depends on that repository and finally the
handler depends on the use case.

## Libraries

Initialization of external libraries are located in `third_party/`

Since `sqlx` is a third party library, it is initialized in `/third_party/database/sqlx.go`

# Cache

The three most significant bottlenecks are 

  1. Input output (I/O) like disk access including database.
  2. Network calls - like calling another API.
  3. Serialization - like serializing or deserializing JSON

We demonstrate how caching results can speed up API response: 

## LRU

To make this work, we introduce another layer that sits between use case and database layer.

`internal/author/repository/cache/lru.go` shows an example of using an LRU cache to tackle the biggest bottleneck. Once we get a result for the first time, we store it by using the requesting URL as its key. Subsequent requests of the same URL will return the result from the cache instead of from the database. 

The request url only exists in the handler layer by accessing it from `*http.Request`.
To pass the request url to our cache layer, we can either pass it down the layers as
a method parameter, or we can use `context` to save this url. In general, `context`
is not the way to store variables but since this url is scoped to this one request,
it is okay.

```go
ctx := context.WithValue(r.Context(), author.CacheURL, r.URL.String())
```

Then in the cache layer, we retrieve it using the same key (`author.CacheURL`)
and we must assert the type to `string`.

```go
url := ctx.Value(author.CacheURL).(string)
```

The code above can panic if the key does not exist. To be safer, we can check if
a value is retrieved successfully.

```go
url, ok := ctx.Value(author.CacheURL).(string)
if !ok {
	// handle if not ok, or call repository layer.
}
```

Using the `url`, we try and retrieve a value from the cache,
```go
val, ok := c.lru.Get(url)
```

If it doesn't exist, we can simply add it to our cache. 
```go
c.lru.Add(url, res)
```

Avoiding I/O bottleneck results in an amazing speed, **11x** more requests/second (328 bytes response size) compared to an already blazing fast endpoint as shown by `wrk` benchmark:

CPU: AMD 3600 3.6Ghz
Storage: SSD

```shell
wrk -t2 -d60 -c200  'http://localhost:3080/api/v1/author?page=1&size=3'
Running 1m test @ http://localhost:3080/api/v1/author?page=1&size=3
  2 threads and 200 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     4.23ms    5.07ms  71.75ms   83.36%
    Req/Sec    40.64k     3.55k   52.91k    68.45%
  4847965 requests in 1.00m, 1.48GB read
Requests/sec:  80775.66
Transfer/sec:     25.27MB
```

Compared to calling database layer:
```shell
wrk -t2 -d60 -c200  'http://localhost:3080/api/v1/author?page=1&size=3'
Running 1m test @ http://localhost:3080/api/v1/author?page=1&size=3
  2 threads and 200 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    70.66ms  116.57ms   1.24s    88.09%
    Req/Sec     3.66k   276.15     4.53k    70.50%
  437285 requests in 1.00m, 136.79MB read
Requests/sec:   7280.82
Transfer/sec:      2.28MB
```

Since a cache stays in the store if it is frequently accessed, invalidating the cache must be done if there are any changes to the stored value in the event of update and deletion. Thus, we need to delete the cache that starts with the base URL of this domain endpoint. 

For example:
```go
func (c *AuthorLRU) Update(ctx context.Context, toAuthor *models.Author) (*models.Author, error) {
	c.invalidate(ctx)

	return c.service.Update(ctx, toAuthor)
}

func (c *AuthorLRU) invalidate(ctx context.Context) {
	url := ctx.Value(author.CacheURL)
	split := strings.Split(url.(string), "/")
	baseURL := strings.Join(split[:4], "/")

	keys := c.lru.Keys()
	for _, key := range keys {
		if strings.HasPrefix(key.(string), baseURL) {
			c.lru.Remove(key)
		}
	}
}
```
## Redis

By using Redis as a cache, you can potentially take advantage of a cluster architecture for more RAM instead of relying on the RAM on current server your API is hosted. Also, the cache won't be cleared like in-memory `LRU` when a new API is deployed.

Similar to LRU implementation above, this Redis layer sits in between use case and database layer.

This Redis library requires payload in a binary format. You may choose the builtin `encoding/json` package or `msgpack` for smaller payload and **7x** higher speed than without a cache. Using `msgpack` over `json` tackles serialization bottleneck.

```go
// marshal 
cacheEntry, err := msgpack.Marshal(res)
// unmarshal
err = msgpack.Unmarshal([]byte(val), &res)
```

```shell
wrk -t2 -d60 -c200  'http://localhost:3080/api/v1/author?page=1&size=3'
Running 1m test @ http://localhost:3080/api/v1/author?page=1&size=3
  2 threads and 200 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     4.05ms    2.56ms  37.48ms   73.63%
    Req/Sec    25.48k     1.45k   30.73k    71.29%
  3039522 requests in 1.00m, 0.93GB read
Requests/sec:  50638.73
Transfer/sec:     15.84MB
````

# Utility

Common tasks like retrieving query parameters or `filters` are done inside `utility` folder. It serves as one place abstract functionalities used across packages.

Note that further packages are in the subdirectories of `utilty`. No files
ar under `utilty` package because it is [unclear from the name](https://go.dev/doc/effective_go#package-names) on what it does. 

# Testing

A testable code is a sign that you have a good code. However, it can be hard to
write not only of a function, but also how different functions work together. 
Going on a tangent, following [SOLID principle](https://en.wikipedia.org/wiki/SOLID)  is a good way to design our code and make it testable. But this repository isn't nearly complex enough to show good examples of each principle. In any case, we shall start with unit tests.

## Unit Testing

Unit testing can be run with

```sh
task test
```
    
Which runs `go test -v ./...`

A quick note, in Go, unit test file is handled by appending `_test` to a file's name. For example, to test `/internal/domain/book/handler/http/handler.go`, we add unit test file by creating `/internal/domain/book/handler/http/handler_test.go`

To perform a unit test we take advantage of go's interface. The layers between
handler, use case and database are loosely-coupled. This is achieved by accepting
an interface and return a struct. This way, you can swap that struct implementation
with something else, say, a mock. So when you run the unit test, it will run
that mock implementation instead of your concrete implementation. This is all
thanks to Go's interface which is implicit.

This repository shows table-driven unit testing strategy  in all three layers.
Both handler and usecase layers swaps the implementation of underneath layer 
with mocks while in repository layer, we use _real_ database in docker to test
against, using `dockertest` library.

### Handler

We explore on how to perform unit testing on creating an Author. There are several things that need to happen, namely:

1. Create a request.
2. Validate request.
3. Call business logic layer underneath it and handle various error that may come up.
    - We are not going to actually call our business logic layer. We use mocks instead.
4. Perform data transformation for user consumption.
5. Compare errors and values.

In general, all unit tests will have `args` and `want` struct. `args` struct is
what we need to supply to the unit test while `want` struct is where we define
what we expect the result is going to be.

Firstly, we create `handler_test.go` file in the same directory. Create a unit
test function called `TestHandler_Create()`.

```go
func TestHandler_Create(t *testing.T) {
	
}
```

In there, we add `CreateRequest` to `args` struct.

```go
type args struct {
    *author.CreateRequest
}
```

In `want` struct, we expect the usecase to return three things, the author an 
error, and http status code. We add an HTTP status response code because our 
handler can return different code depending on the result.

What is strange here is we included a `usecase` struct. This is the struct that
contains values that our mock will return.

```go
type want struct {
	usecase struct {
        *author.Schema
        error
    }
    response *author.GetResponse
    status int
    error
}
```

The final struct embeds both structs, and we give a name to it. Naming the test
makes it easier to run individual tests.

```go
type test struct {
    name string
    args
    want
}
```

Now that we have all necessary structs, we can begin with our table-driven tests.
It is just a matter of filling `test` struct with our values.

```go
tests := []test{
	{
        name: "simple",
        args: args {
            CreateRequest: &author.CreateRequest{
                FirstName:  "First",			
                MiddleName: "Middle",			
                LastName:   "Last",
	            Books:      nil,
            }   		
        },
        want: want{
				usecase: struct {
					*author.Schema
					error
				}{
					&author.Schema{
						ID:         1,
						FirstName:  "First",
						MiddleName: "Middle",
						LastName:   "Last",
						CreatedAt:  time.Now(),
						UpdatedAt:  time.Now(),
						DeletedAt:  nil,
						Books:      make([]*book.Schema, 0),
					},
					nil,
				},
				response: &author.GetResponse{
					ID:         1,
					FirstName:  "First",
					MiddleName: "Middle",
					LastName:   "Last",
					Books:      make([]*book.Schema, 0),
				},
				err:    nil,
				status: http.StatusCreated,
			},
}
```

As you can see, the _author_ record gets replicated three times. But there is a
subtle difference between the input and output. In the `args` struct, we only
supply three things, the names, but left the books as empty. In the `response`
struct, it has the names, but also a zero-length array of books! This is an
intended behaviour - I purposely wanted to return an array no matter what.

So within this `tests` slice, we can add as many unit tests as we want. We could
try giving invalid input, input that fails validation, etc.  

To run the tests, we loop over this slice of tests:

```go
for _, test := range tests {
    t.Run(test.name, func(t *testing.T) {
        
    }
}
```

We use `httptest` package to call our tests by creating a `writer`(request) to 
call the handler, and a `recorder` to receive response from the handler.

```go
rr := httptest.NewRequest(http.MethodPost, "/api/v1/author", <body>)
ww := httptest.NewRecorder()
```

The request points to the URL of the endpoint, and we make a `POST` request to it.
Since we are sending a JSON payload, we send it in the third argument. It accepts
an `io.Reader` so we need to encode our JSON payload into `bytes.Buffer`:

```go
var buf bytes.Buffer
err = json.NewEncoder(&buf).Encode(test.args.CreateRequest)

rr := httptest.NewRequest(http.MethodPost, "/api/v1/author", &buf)
```

This is a good place to assert that no error has happened.

```go
err = json.NewEncoder(&buf).Encode(test.args.CreateRequest)
assert.Nil(t, err)
```

To call our handler, we need to instantiate it. It is created from `RegisterHTTPEndPoints()`.

```go
h := RegisterHTTPEndPoints(router, val, uc)
```

This function requires three dependencies. The `router` and `validator` are easy:

```go
router := chi.NewRouter()
val := validator.New()
```

The final dependency requires a bit of work. The handler depends on the usecase
interface, and it in turn calls the appropriate concrete implementation. For our
unit test, we swap out the implementation with a mock. And this mock returns
value from our `want.usecase` struct. Now our unit test can work in isolation, 
and do not depend on any layer underneath it!

Create a new file called `usecase_mock.go`. Declare a new mock struct and within it,
contains a field that matches our usecase signature by looking at the usecase 
interface.

`usecase.go`
```go
type UseCase interface {
    Create(ctx context.Context, a *author.CreateRequest) (*author.Schema, error)
}
```

`usecase_mock.go`
```go
type AuthorUseCaseMock struct {
    CreateFunc func(ctx context.Context, a *author.CreateRequest) (*author.Schema, error)
}
```
Notice that we append the `Create()` method with `Func` field. Now that we have
the struct defined, we add a concrete implementation from it.

`usecase_mock.go`
```go
func (m *AuthorMock) Create(ctx context.Context, a *author.CreateRequest) (*author.Schema, error) {
	return a.CreateFunc(ctx, req)
}
```

Now that we have a usecase mock, we can now give the `uc` variable for `RegisterHTTPEndPoints()` function. Using `AuthorUseCaseMock` struct from `mock` package, we initialize `CreateFunc` field from it. Then, it is just a matter of returning the values to what we have defined in our `want` struct.

`handler_test.go`
```go
uc := &mock.AuthorUseCaseMock{
    CreateFunc: func(ctx context.Context, a *author.CreateRequest) (*author.Schema, error) {
        return test.want.usecase.Schema, test.want.usecase.error
    },
}
```

We finally have all of our dependencies initialized. Now we can call `Create()`
method. We pass in the writer(`ww`) and a recorder `rr` into it - which matches our
handler signature (`Create(w http.ResponseWriter, r *http.Request)`)

```go
h := RegisterHTTPEndPoints(router, val, uc)
h.Create(ww, rr)
```

Response is recorded into `ww` variable. To receive the response, we decode from
`ww.Body` into `author.Schema` struct:

```go
var got author.GetResponse
if err = json.NewDecoder(ww.Body).Decode(&got); err != nil {
    t.Fatal(err)
}
```

Finally, we can do some assertions to check if the returned response matches with
what we expect.

```go
assert.Equal(t, ww.Code, test.status)
assert.Equal(t, &got, test.want.response)
```

While `go test ./...` runs all tests, we can choose to run only this specific test. We `cd` into the directory and use `-run` to specify the <function name/test name>. `-run` can also accept regex

```shell
cd internal/domain/author/handler
go test -run="TestHandler_Create/simple"

PASS
ok      github.com/gmhafiz/go8/internal/domain/author/handler   0.010s
```

There are a lot of things going on in this unit test. It is very verbose, but it
is clear on what happens here. A table-test allows us to quickly construct arguments,
wants and what we expect. Constructing a mock file can be tedious, so a tool like
[mirip](https://github.com/gmhafiz/mirip) can be used to generate a mock from your
interface.

### Use Case

The idea is the same as unit testing a handler. We have a set of arguments, what
is expected from it, and a slice of `test` struct that we iterate to test against.

This time, we do not have to worry about write and recorder. We only need to 
instantiate usecase along with its dependencies. To make this simple, we will
only mock database(repository struct).

The `Create()` method expects a `context` and `*author.CreateRequest` and returns
`*author.Schema` and an `error`.
```go
type args struct {
    *author.CreateRequest
}
type want struct {
    *author.Schema
    error
}
```

Our `test` struct becomes

```go
type test struct {
    name string
    args
    want
}
```

Like handler unit tests, we fill in the `test` slice with our data

```go
tests := []test{
    {
        name: "simple",
        args: args{
            CreateRequest: &author.CreateRequest{
                FirstName:  "First",
                MiddleName: "Middle",
                LastName:   "Last",
                Books:      nil,
            },
        },
        want: want{
            Author: &author.Schema{
                ID:         1,
                FirstName:  "First",
                MiddleName: "Middle",
                LastName:   "Last",
                CreatedAt:  time.Time{},
                UpdatedAt:  time.Time{},
                DeletedAt: nil,
                Books:     nil,
            },
            err: nil,
        },
    },
}
```

To instantiate a usecase, we call the `New()` function.

```go
uc := New(nil, repoAuthor, nil, nil, nil, nil)
```

We only care about CRUD at this stage, so we only need to mock out the repository
layer. We start by creating a mock file and create a struct containing methods that 
matches the signature defined in repository interface.

`postgres.go`
```go
type Repository interface {
    Create(ctx context.Context, a *author.CreateRequest) (*author.Schema, error)
}
```
`postgres_mock.go`
```go
package database

type AuthorMock struct {
	CreateFunc func(ctx context.Context, a *author.CreateRequest) (*author.Schema, error)
}
```

Then we implement `CreateFunc` method.

```go
package database

type AuthorMock struct {
    CreateFunc func(ctx context.Context, a *author.CreateRequest) (*author.Schema, error)
}

func (m *AuthorMock) Create(ctx context.Context, a *author.CreateRequest) (*author.Schema, error) {
	return m.CreateFunc(ctx, a)
}
```

As mentioned before, we can use [mirip](https://github.com/gmhafiz/mirip) to
automatically generate this mock file.

```go
//go:generate mirip -rm -out postgres_mock.go . Repository
type Repository interface {...
```

Now that we have repository mock, we can start looping through `tests` variable.
```go
for _, test := range tests {
    t.Run(test.name, func(t *testing.T) {

    }
}
```

Inside, we declare `&repository.AuthorMock` for repository mock. It returns
the author and error that we want as declared in the table-test.

```go
repoAuthor := &repository.AuthorMock{
    CreateFunc: func(ctx context.Context, r *author.CreateRequest) (*author.Schema, error) {
        return test.want.Author, test.want.err
    },
}

uc := New(nil, repoAuthor, nil, nil, nil)
```

With the usecase declared, we can call its `Create()` method.

```go
got, err := uc.Create(context.Background(), test.args.CreateRequest)
```

Finally, we perform a couple of assertions to check for error and response

```go
assert.Equal(t, test.want.error, err)
assert.Equal(t, test.want.Author, got)
```

Run the test with

```shell
cd internal/domain/author/usecase
go test -run="TestAuthorUseCase_Create/simple"

PASS
ok      github.com/gmhafiz/go8/internal/domain/author/usecase   0.004s
```

### Repository

Unit testing repository layer is different from above in the way that we test
them against real database using Docker, instead of using mocks. It is a lot
more complex to set up because now we need to do at least two things:

1. Instantiate a new database in Docker
2. Perform migration to create the tables
3. Seed, if necessary

A question arise wondering - if we are testing against a real database, even if
it is a temporary one, is it really a unit test? Should it not be called an 
integration test? Honestly it does not matter what it is called now. The fact
that we are still in one layer, the repository layer, not touching the usecase,
or handler, I'd consider that we are still doing unit tests. Moreover, I find
that doing mocks for database is flaky and doesn't capture many edge cases, so
potentially rendering them useless.

Moving on, to set up, we use `TestMain`. It will run before all unit tests in this package. The code is basically a copy-paste from https://github.com/ory/dockertest. We 
can customize the image (`postgres`) and version using tag (`15`). The username,
password and database name is not important, and they can be anything. These 
databases will be automatically shut down. In spite of spinning a database for 
these tests, running these unit tests are still quick. For example, running all
15 CRUD tests in this `database` package takes only **2** seconds.

```go
func TestMain(m *testing.M) {
	// 1. the place where we create our database container that starts really fast
	// 2. Run migration
}

```
Once the database is up, we need to create initial tables. Thus, we call
`migrate("up")`. Full code is in `postgres_test.go`.

Now that we have database set up, we write our first repository unit test on
`Create()`.

```go
func TestAuthorRepository_Create(t *testing.T) {

}
```

Like other unit tests, we supply `args`, `want`, and `test` struct. Look at
`type Repository interface` to know the function's signature to infer what are
needed and what it returns.

```go
type args struct {
    author *author.CreateRequest
}

type want struct {
    author *author.Schema
    err    error
}

type test struct {
    name string
    args
    want
}
```

In the `test` slice, we supply the values. In this particular example, a simple 
one. You may browse the codebase for more test scenarios.

```go
startTime := time.Now()

tests := []test{
    {
        name: "normal",
        args: args{
            author: &author.CreateRequest{
                FirstName:  "First",
                MiddleName: "Middle",
                LastName:   "Last",
                Books:      nil,
            },
        },
        want: want{
            author: &author.Schema{
                    ID:         1,
                    FirstName:  "First",
                    MiddleName: "Middle",
                    LastName:   "Last",
                    CreatedAt:  time.Now(),
                    UpdatedAt:  time.Now(),
                    DeletedAt:  nil,
                },
            err: nil,
        },
    },
}
```

Eventually we need to call `Create()` method from this repository. For that we 
need to instantiate a repository, and this repository requires a database client.

```go
var (
    dockerDB *db
)

type dockerDB struct {
    Conn *sql.DB
    Ent  *gen.Client
}

func dbClient() *gen.Client {
    sqlxDB := sqlx.NewDb(DockerDB.Conn, "postgres")
    drv := entsql.OpenDB(dialect.Postgres, sqlxDB.DB)
    client := gen.NewClient(gen.Driver(drv))
    DockerDB.Ent = client
    
    return client
}

client := dbClient()
```
This database client is created from `sqlx`, and since we are using `ent` ORM, we create a client from it too. Then we store the `client` in a variable local to `database` package. This way, the database client is accessible to all other unit tests.

With a database client, we can now create a repository

```go
repo := New(client)
```

To run our table-test, like before, we iterate the `test` slice

```go
for _, test := range tests {
    t.Run(test.name, func (t *testing.T) {
        ctx := context.Background()
        
        created, err := repo.Create(ctx, test.args.author)
    })
}
```
It returns two values, `created` and `err`. So we assert them. Since we used real database, it returns values that we cannot know in advance like `created_at` field. If you scroll up a bit, you will notice that we declared a `startTime`. Using this trick I learned in the [#gopher](https://gophers.slack.com) Slack channel, you simply assert that both `created_at`, and `updated_at` timestamp values must be between that `startTime` value and the same as the generated timestamps.

```go
assert.Equal(t, err, test.want.err)
assert.Equal(t, created.ID, test.want.author.ID)
assert.Equal(t, created.FirstName, test.want.author.FirstName)
assert.Equal(t, created.MiddleName, test.want.author.MiddleName)
assert.Equal(t, created.LastName, test.want.author.LastName)

assert.True(t, startTime.Before(created.CreatedAt) || startTime.Equal(created.CreatedAt))
assert.True(t, startTime.Before(created.UpdatedAt) || startTime.Equal(created.UpdatedAt))
assert.Equal(t, test.want.author.DeletedAt, created.DeletedAt)
```

To run

```shell
cd internal/domain/author/repository/database
go test -run="TestAuthorRepository_Create/normal"

PASS
ok      github.com/gmhafiz/go8/internal/domain/author/repository/database       2.320s

```

To get more test coverage, more tests need to be added to the table-test. For 
example, inserting empty payload, test for errors, and inserting author with
attached books. 

When doing multiple inserts, the ID will be increased. Remember that we are 
testing against real database. So your expected ID should follow the increment.

The way tests are laid out, there is a single database client used by all
unit tests. Alternatively, you can choose to have one database client for each
of `Create()`, `Read()`, `Update()` and `Delete()` which would have mean 4 
separate databases in its own Docker container.

Note that in doing a `Read()` unit test, if you choose to seed by doing an 
insert, instead of using SQL script before reading, you are already doing an 
integration test. An integration test is simply several unit tests that work 
together.

In conclusion, unit testing repository layer is more verbose as it needed a
third party library. However, the structure is still similar, with the addition
of setting up the database. Another advantage is you can inspect the values
inside the database by looking at the `databaseUrl` variable inside `TestMain()`.

## End-to-End Test

TODO: E2E tests are still in progress.

Technically End-to-End test (e2e test) can be done separately in another program and language. Having e2e binary integrated in the project has the advantage of reusing structs and migration which will be explained down below. 

The idea here is to run our application isolated in a container (along with database) and the e2e program calls known API of this program and checks if the output is what is expected.  

It creates two docker containers, one for the API, and second for postgres. Once the containers have started. It runs the app, and then the e2e binary.

Remember the `Server` struct in `internal/server/server.go` file? The `New()` function is called both by our API and e2e binary. We can also call `Migrate()` function because the e2e test uses the same `Server` struct as our API.

In our actual e2e implementation use cases, we can perform various CRUD operations.

For example, in an empty database, we expect no books should be returned.

```go
func testEmptyBook(t *E2eTest) {
	// call our API endpoint
	resp, err := http.Get(fmt.Sprintf("http://localhost:%s/api/v1/books", t.server.Config().Api.Port))
	
	// The return should be an empty array
    if !bytes.Equal(expected, got) {
        log.Printf("handler returned unexpected body: got %v want %v", string(got), expected)
    }
}
```

### Run e2e test

Start

    task dockertest

or

```shell
 cd docker-test && docker-compose down -v --build && docker-compose up -d
 docker exec -t go8_container_test "/home/appuser/app/e2e"
```

# TODO

 - [ ] Fix end-to-end test
 - [ ] Complete HTTP integration test
 - [x] Better return response
 - [x] LRU cache
 - [X] Redis Cache
 - [ ] Tracing
 - [ ] Metric

# Acknowledgements

 * https://quii.gitbook.io/learn-go-with-tests/questions-and-answers/http-handlers-revisited
 * https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html
 * https://github.com/moemoe89/integration-test-golang
 * https://github.com/george-e-shaw-iv/integration-tests-example
 * https://gist.github.com/Oxyrus/b63f51929d687c1e20cda3447f834147
 * https://github.com/sno6/gosane
 * https://github.com/AleksK1NG/Go-Clean-Architecture-REST-API
 * https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1
 * https://github.com/arielizuardi/golang-backend-blog
 
# Appendix

## Dev Environment Installation

For Ubuntu:

```shell
sudo apt update && sudo apt install git curl build-essential jq
wget https://go.dev/dl/go1.19.4.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.19.4.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
echo 'PATH=$PATH:/usr/local/go/bin' >> ~/.bash_aliases
echo 'PATH=$PATH:$HOME/go/bin' >> ~/.bash_aliases
source ~/.bashrc
go install golang.org/x/tools/...@latest

curl -s https://get.docker.com | sudo bash
sudo usermod -aG docker ${USER}
newgrp docker
su - ${USER} # or logout and login

sudo curl -L "https://github.com/docker/compose/releases/download/v2.9.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```
