# Introduction

A starter kit for Go API development. Heavily based on [goapp](https://github.com/bnkamalesh/goapp)
that does an excellent job of organizing things together as well inspired by [How I write HTTP
 services after eight years](https://pace.dev/blog/2018/05/09/how-I-write-http-services-after-eight-years.html).
 However I wanted to use [chi router](https://github.com/go-chi/chi) and [sqlboiler](https://github.com/volatiletech/sqlboiler/)
 which are more common in the Go community.

Just like the project it is based on, this kit tries to follow the [Standard Go Project Layout
](https://github.com/golang-standards/project-layout) to make project structure familiar to a Go
 developer.

It is still in early stages and I do not consider it is completed until all integration test and
 input validation are done.

In short, this kit is a Go + Postgres + Chi Router + SqlBoiler starter kit for API development.

# Motivation

On the topic of API development, the Go community is split between recommending a framework (like
 [echo](https://github.com/labstack/echo), [gin](https://github.com/gin-gonic/gin), 
   [buffalo](http://gobuffalo.io/) and using standard Go library plus other libraries you need. 
   However, starting small and sticking to built in `net/http` library plus a few other
    well known libraries give a lot more flexibility to designing an API since you can always
     plug and play any functionality you want. 

# Features

This kit is composed of standard Go library together with well known libraries to
 manage things like router, database query and migration support. Technically it supports 
 [other databases](https://github.com/volatiletech/sqlboiler#supported-databases) as well. 

  - [Chi Router](https://github.com/go-chi/chi) 
  - [Sqlboiler ORM](https://github.com/volatiletech/sqlboiler/)
  - Database migration with [golang-migrate](https://github.com/golang-migrate/migrate/)
  - Cache result with [Redis](https://redis.io) using [msgpack](https://msgpack.org) 
  - Input [validation](https://github.com/go-playground/validator) that return multiple error
   strings
  - Scans and auto-generate [Swagger](https://github.com/swaggo/swag) docs using a declarative
   comments format 
  - Request log that logs each user uniquely based on host address
  - Cors
  - HTTP Integration Test
  - Pagination
  - Yaml file for configuration

It has few dependencies and replacing one library to another is easy as long as it adheres to
 standard Go library interface.


# Getting It

    git clone https://github.com/gmhafiz/go8
    cd go8

# Setup

A. Have both a postgres database and a redis instance ready.

If not, you can run the following command if you have `docker-compose` installed:
 
    docker-compose up -d postgres redis

B. This project uses [Task](https://github.com/go-task/task) to handle various tasks such as
 migration, generate swagger docs, build and run the app. It is essentially a [sh interpreter
 ](https://github.com/mvdan/sh). Only requirement is to download the binary and append to your `PATH` variable.
  - Install task runner binary bash script:

    
    scripts/install-task.sh

  - And put this binary in your path if not exists
  
    
    echo 'PATH=$PATH:$HOME/.local/bin' >> ~/.bashrc
    source ~/.bashrc        

`Tasl` tasks are defined inside `Taskfile.yml` file. A list of tasks available can be viewed with:
                                                     
    task -l   # or
    task list # which maps to `task -l`

Once `Task` is installed, setup can be
 initiated by the
 following
 command:

    task init
    
This copies example configurations for the app, `sqlboiler` and `Task` to its respective .yml
 files as well as syncs dependencies
Then open the files at `config/dev.yml`, `sqlboiler.toml`, `.env` and fill in your own configurations


C. This project uses [golang-migrate](https://github.com/golang-migrate/migrate/) to handle
 database migrations and [Sqlboiler ORM](https://github.com/volatiletech/sqlboiler/) to handle
  database queries. These tools can be installed with:
  

    task install-tools
  

## Migration

Migration is a good step towards having a versioned database and makes publishing to a production
 server a safe process.
    
### Create Migration

Using `Task`, creating a migration file is done by the following command. Name the file after
 `NAME=`. 

    task migrate-create NAME=create_a_table

### Migrate up

After you are satisfied with your `.sql` files, run the following command to migrate your database.

    task migrate
    

## Database Generate Models and ORMs

SqlBoiler treats your database as source of truth. It connects to your database, read its schema
 and generate appropriate models and query builder helpers written in Go. Utilizing a type-safe
  query building allows compile-time error checks. 

Generate ORM with:    
    
    task gen-orm

Generated files are as defined in the `sqlboiler.toml` file. This command needs to be run after
 every migration changes are done.

# Test

Install testify testing framework with

    go get github.com/stretchr/testify

# Run

## Local

Conventionally, all apps are placed inside the `cmd` folder.

Using `Task`:

    task run

Without `Task`
    
    go run cmd/go8/go8.go 
    
## Docker

You can build a docker image with the app with its config files. Docker needs to be installed
 beforehand.

     task docker-build

Run the following command to build a container from this image. `--net=host` tells the container
 to use local's network so that it can access local's database.

    task docker-run

In terminal you will see the API runs at port 3080 and a log of available paths

    2020-09-25T09:43:37+10:00 INF internal/server/http/http.go:40 > starting at :3080 service=go8
    2020-09-25T09:43:37+10:00 INF internal/server/http/routes.go:50 >  routes={"method":"POST","path":"/api/v1/author"} service=go8
    2020-09-25T09:43:37+10:00 INF internal/server/http/routes.go:50 >  routes={"method":"GET","path":"/api/v1/author/{authorID}"} service=go8
    2020-09-25T09:43:37+10:00 INF internal/server/http/routes.go:50 >  routes={"method":"GET","path":"/api/v1/authors"} service=go8
    2020-09-25T09:43:37+10:00 INF internal/server/http/routes.go:50 >  routes={"method":"POST","path":"/api/v1/book"} service=go8
    2020-09-25T09:43:37+10:00 INF internal/server/http/routes.go:50 >  routes={"method":"GET","path":"/api/v1/book/{bookID}"} service=go8
    2020-09-25T09:43:37+10:00 INF internal/server/http/routes.go:50 >  routes={"method":"DELETE","path":"/api/v1/book/{bookID}"} service=go8
    2020-09-25T09:43:37+10:00 INF internal/server/http/routes.go:50 >  routes={"method":"GET","path":"/api/v1/books"} service=go8
    2020-09-25T09:43:37+10:00 INF internal/server/http/routes.go:50 >  routes={"method":"GET","path":"/health/liveness"} service=go8
    2020-09-25T09:43:37+10:00 INF internal/server/http/routes.go:50 >  routes={"method":"GET","path":"/health/readiness"} service=go8
    2020-09-25T09:43:37+10:00 INF internal/server/http/routes.go:50 >  routes={"method":"GET","path":"/*"} service=go8


## Docker Compose

If you have `docker-compose` installed, you may run the app with the following command. Docker
-compose binary must be installed beforehand. 

    docker-compose up -d

Both Postgres and redis ports are mapped to local machine. To allow `api` container to reach the
 database and redis.


# Swagger docs
     
Edit `cmd/go8/go8.go` `main()` function host and BasePath  

    // @host localhost:3080
    // @BasePath /api/v1

   
Generate with

    task swagger
    
Access at

    http://localhost:3080/swagger

The command `swag init` scans the whole directory and looks for [swagger's declarative comments](https://github.com/swaggo/swag#declarative-comments-format)
 format.

Custom theme is obtained from [https://github.com/ostranme/swagger-ui-themes](https://github.com/ostranme/swagger-ui-themes)


# Cache

Redis cache is by default 5 seconds. It is set by the `Set()` method in `store.go` file.

# Tooling

Various tooling are included within the `Task` runner

  * `task fmt`
    * Runs `go fmt ./...` to lint Go code
  * `task tidy`
    * Runs `go mod tidy` to sync dependencies
  * `task vet`
    * Quickly catches compile error
  * `task golint`
    * Runs an opinionated code linter from https://golangci-lint.run/

# Structure
    
1. Entry point is at `cmd/go8/go8./go`
2. Api Routes are defined by `chi` router in `internal/server/http/routes.go`
3. Handlers are defined under `internal/server/http` folder
4. Each entity (`book` and `author`) is in their own microservice in `internal/domain` folders. 
This makes the layout confusing but allows dependency injection for integration testing purpose.
5. Migration `.sql` files goes under `database/migrations` folder.
6. Config `.yml` files goes under `config` folder. You can place `dev.yml`, `test.yml`, `prod.yml` 
under this folder. Note: all `.yml` and `.toml` files are ignored by version control.

# TODO

 - Complete HTTP integration test
 - use [xID](https://github.com/rs/xid) for table ID primary key
