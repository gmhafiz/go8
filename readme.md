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
  - Request log
  - Cors
  - HTTP Integration Test
  - Pagination
  - Yaml file for configuration

It has few dependencies and replacing one library to another is easy as long as it adheres to
 standard Go library interface.

# Setup

  - Have an empty Postgres database ready
  - Copy configuration files and fill in database and api details 
    - `cp config/dev.yml.example config/dev.yml`
    - `cp sqlboiler.toml.example sqlboiler.toml`
  - Install the following tools. Instructions are in the next sections.
    - [golang-migrate](https://github.com/golang-migrate/migrate/)
    - [Sqlboiler ORM](https://github.com/volatiletech/sqlboiler/)

## Migration

Migration is a good step towards having a versioned database and makes publishing to a production
 server a safe process.
 
 While there are many ways to [install](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
 `golang-migrate`, simplest way to get migration going is to download its binary. Latest releases
  are at its [releases page](https://github.com/golang-migrate/migrate/releases).

Download binary

    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.11.0/migrate.linux-amd64.tar.gz | tar xvz

Add the binary to your $PATH

    mkdir -p ~/.local/bin
    mv migrate.linux-amd64 ~/.local/bin/migrate
    source ~/.bashrc
    
### Create Migration

    migrate create -ext sql -dir database/migrations -format unix create_books_table
    migrate create -ext sql -dir database/migrations -format unix create_authors_table
    migrate create -ext sql -dir database/migrations -format unix create_book_authors_table


### Migrate up

    migrate -database "postgres://127.0.0.1/db?sslmode=disable&user=user&password=pass" -path database/migrations up
    

## Database Generate Models and ORMs

SqlBoiler treats your database as source of truth. It connects to your database, read its schema
 and generate appropriate models and query builder helpers written in Go. Utilizing a type-safe
  query building allows runtime error checks. 

### Install

    GO111MODULE=off go get -u -t github.com/volatiletech/sqlboiler
    GO111MODULE=off go get github.com/volatiletech/sqlboiler/drivers/sqlboiler-psql
         
Fill in settings in `sqlboiler.toml` file

Generate ORM with:    
    
    sqlboiler --wipe --add-soft-deletes psql

Generated files are as defined in the `sqlboiler.toml` file. This command needs to be run after
 every migration changes are done.

# Test

Install testify testing framework with

    go get github.com/stretchr/testify

# Run

## Local

Conventionally, all apps are placed inside the `cmd` folder.

    go run cmd/go8/go8.go 
    
## Docker

You can build a docker image with the app with its config files.

     docker build -t go8 -f docker/Dockerfile .

Run the following command to build a container from this image. `--net=host` tells the container
 to use local's network so that it can access local's database.

    docker run -p 3080:3080 --rm -it --net=host go8


## Docker Compose

If you have `docker-compose` installed, you may run the app with the following command. 

    docker-compose up -d

Both Postgres and redis ports are mapped to local machine. To allow `api` container to reach the
 database and redis.



# Structure
    
1. Entry point is at `cmd/go8/go8./go`
2. Api Routes are defined by `chi` router in `internal/server/http/routes.go`
3. Handlers are defined under `internal/server/http` folder
4. Each entity (`book` and `author`) is in their own microservice in `internal/service` folders. 
This makes the layout confusing but allows dependency injection for integration testing purpose.
5. Migration `.sql` files goes under `database/migrations` folder.
6. Config `.yml` files goes under `config` folder. You can place `dev.yml`, `test.yml`, `prod.yml` 
under this folder. Note: all `.yml` and `.toml` files are ignored by version control.

# TODO

 - Complete HTTP integration test
 - Use sqlboiler as a library and make an executable under folder `cmd/sqlboiler` to have a single
  `yml` config file.
 - Use golang-migrate as a library and make an executable under folder `cmd/migrate`
 - Swagger documentation
 - use [xID](https://github.com/rs/xid) for table ID primary key
 - consider using [mage](https://github.com/magefile/mage) to simplify build process

