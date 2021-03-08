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

However, I wanted to use [chi router](https://github.com/go-chi/chi) which is more common in the community, [sqlx](https://github.com/jmoiron/sqlx) for database operations and design towards more like [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html).

This kit tries to follow the [Standard Go Project Layout](https://github.com/golang-standards/project-layout) to make project structure familiar to a Go developer.

It is still in early stages, and I do not consider it is completed until all integration tests are completed.

In short, this kit is a Go + Postgres + Chi Router + sqlx + unit testing starter kit for API development.

# Motivation

On the topic of API development, there are two opposing camps between a framework (like [echo](https://github.com/labstack/echo), [gin](https://github.com/gin-gonic/gin), [buffalo](http://gobuffalo.io/) and starting small and only add features you need. 

However , starting small and adding  features aren't that straightforward. Also, you will want to structure your project in such a way that there are clear separation of functionalities for your controller, business logic and database operations. 

This is the idea behind [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html). This way, it is easy to switch whichever library to another of your choice.


# Features

This kit is composed of standard Go library together with some well-known libraries to manage things like router, database query and migration support.

  - [x] Router/Mux with [Chi Router](https://github.com/go-chi/chi)
  - [x] Database Operations with [sqlx](https://github.com/jmoiron/sqlx)
  - [x] Database migration with [golang-migrate](https://github.com/golang-migrate/migrate/)
  - [x] Input [validation](https://github.com/go-playground/validator) that return multiple error strings
  - [x] Read all configurations using a single `.env` file
  - [x] (optional) Request log that logs each user uniquely based on host address
  - [x] Cors
  - [x] Scans and auto-generate [Swagger](https://github.com/swaggo/swag) docs using a declarative comments format 
  - [x] Custom model JSON output
  - [x] Filters
  - [x] Uses [Task](https://taskfile.dev) to simplify various tasks 
  - [x] Unit testing of repository, use case, and handler
  - [x] End-to-end test using ephemeral docker containers
  

It has few dependencies and replacing one library to another is easy.


# Quick Start

You need to [have a go installation](#appendix) (>= v1.13) and put into path as well as [git](#appendix). Optionally `docker` and `docker-compose` for easier start up.

Get it

    git clone https://github.com/gmhafiz/go8
    cd go8

Fill in your database credentials in `.env` by making a copy of `env.example` first.

    cp env.example .env

Have a database ready either by installing them yourself or the following command. the `docker-compose.yml` will use database credentials set in `.env` file which is initialized by the previous step. Optionally, you may want redis as well.

    docker-compose up -d postgres

Once the database is up you may run the migration with,

    go run cmd/extmigrate up

Run the API with

    go run cmd/go8/main.go


You will see the address the API is running at as well as all registered routes.

    2021/01/26 18:45:22 serving at 0.0.0.0:3080
    2021/01/26 18:45:22 path: /api/v1/books/ method: GET 
    2021/01/26 18:45:22 path: /api/v1/books/ method: POST 
    2021/01/26 18:45:22 path: /api/v1/books/{bookID} method: GET 
    2021/01/26 18:45:22 path: /api/v1/books/{bookID} method: PUT 
    2021/01/26 18:45:22 path: /api/v1/books/{bookID} method: DELETE 
    2021/01/26 18:45:22 path: /health/liveness method: GET 
    2021/01/26 18:45:22 path: /health/readiness method: GET 


To use, follow examples in the `examples/` folder

    curl --location --request GET 'http://localhost:3080/api/v1/books'


# Tools

While the above quick start is sufficient to start the API, some tools are included for easier task management.

A. This project uses [Task](https://github.com/go-task/task) to handle various tasks such as migration, generation of swagger docs, build and run the app. It is essentially a [sh interpreter](https://github.com/mvdan/sh).

Install task runner binary bash script:

    sudo ./scripts/install-task.sh

This installs `task` to `/usr/local/bin/task` so `sudo` is needed.

`Task` tasks are defined inside `Taskfile.yml` file. A list of tasks available can be viewed with:

    task -l   # or
    task list

## Migration

Migration is a good step towards having a versioned database and makes publishing to a production server a safe process.

### Create Migration

Using `Task`, creating a migration file is done by the following command. Name the file after `NAME=`.

    task migrate-create NAME=create_a_tablename

Write your schema in pure sql in the 'up' version and any reversal in the 'down' version of the files.
 
### Migrate up

After you are satisfied with your `.sql` files, run the following command to migrate your database.

    task migrate

To migrate one step

    task migrate-step n=1
      
### Rollback
    
To roll back migration

    task migrate-rollback n=1

Further `golang-migrate` commands are available in its [documentation (postgres)](https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md)

# Run

## Local

Conventionally, all apps are placed inside the `cmd` folder.

If you have `Task` installed, the server can be run with:

    task run

or without `Task`, just like in quick start section:

    go run cmd/go8/main.go

## Docker

You can build a docker image with the app with its config files. Docker needs to be installed beforehand.

     task docker-build

Run the following command to build a container from this image. `--net=host` tells the container to use local's network so that it can access host database.

    docker-compose up -d postgres # If you haven't run this from quick start 
    task docker-run

### docker-compose

If you prefer to use docker-compose instead, both server and the database can be run with:

    task docker-compose-start

# Tooling

Various tooling are included within the `Task` runner

  * `task fmt`
    * Runs `go fmt ./...` to lint Go code
    * `go fmt` is part of official Go toolchain that formats your code into an opinionated format.
  * `task tidy`
    * Runs `go mod tidy` to sync dependencies.
  * `task vet`
    * Quickly catches compile error.
  * `task golint`
    * Runs an opinionated code linter from https://golangci-lint.run/
  * `task security`
    * Runs opinionated security checks from https://github.com/securego/gosec

All of these can be run with:

    task check

# Swagger docs

Swagger UI allows you to play with the API from a browser

![swagger UI](assets/swagger.png)
     
Edit `cmd/go8/go8.go` `main()` function host and BasePath  

    // @host localhost:3080
    // @BasePath /api/v1

   
Generate with

    task swagger
    
Access at

    http://localhost:3080

The command `swag init` scans the whole directory and looks for [swagger's declarative comments](https://github.com/swaggo/swag#declarative-comments-format) format.

Custom theme is obtained from [https://github.com/ostranme/swagger-ui-themes](https://github.com/ostranme/swagger-ui-themes)

# Structure

This project mostly follows the structure documented at [Standard Go Project Layout](https://github.com/golang-standards/project-layout).

In addition, this project also tries to follow [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) where each functionality are separated into different files.

## Starting Point
Starting point of project is at `cmd/go8/main.go`

![main](assets/main.png)

`s.Init()` in `internal/server/server.go` simply creates a new server, initializes server configuration, database, router, and domains. Lastly `s.Run()` starts the server.

![init](assets/init.png)


## Configurations
![configs](assets/configs.png)

All environment variables are read into specific `Configs` struct initialized in `configs/configs.go`.Each of the embedded struct are defined in its own file of the same package where its fields are read from either environment variable or `.env` file.

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

   ```
    touch configs/elasticsearch.go
   ```
    
2. Create a new struct for your type

    ```
    type Elasticsearch struct {
        Address  string
        User     string
        Password string
    }
    ```
    
3. Add a constructor for it

    ```
    func ElasticSearch() *Elasticsearch {
        return &Elasticsearch{
            Address:  os.Getenv("ELASTICSEARCH_ADDRESS"),
            User:     os.Getenv("ELASTICSEARCH_USER"),
            Password: os.Getenv("ELASTICSEARCH_PASS"),
        }
    }
    ``` 

4. Add to `.env` of the new environment variables

    ```
    ELASTICSEARCH_ADDRESS=http://localhost:9200
    ELASTICSEARCH_USER=user
    ELASTICSEARCH_PASS=password
    ```

Limiting the number of connection pool avoids ['time-slicing' of the CPU](https://github.com/brettwooldridge/HikariCP/wiki/About-Pool-Sizing). Use the following formula to determine a suitable number
 
    number of connections = ((core_count * 2) + effective_spindle_count)    

## Database

Migrations files are stored in `database/migrations` folder. [golang-migrate](https://github.com/golang-migrate/migrate) library is used to perform migration using `task` commands.

## Router

Router or mux is created for use by `Domain`.

Middleware that affects all routes such as CORS, request log and panic recoverer can be registered here.

## Domain

Let us look at how this repository attempts at Clean Architecture. A domain consists of: 

  1. Handler (Controllers)
  2. Use case (Use Cases)
  3. Repository (Entities)
  
Let us look at how `repository` is implemented.

![clean architecture](assets/CleanArchitecture.jpeg)

### Repository

Starting with inner most circle, `Entities`. This is where all database operations are handled. Looking at the `internal/domain/health` folder:

![book-domain](assets/domain-health.png)

Interfaces for both use case and repository are on its own file under the `health` package while its implementation in its separate `usecase` and `repository` package respectively.

Starting with `internal/domain/health/repository.go`

    type Repository interface {
    	Readiness() error
    }

is implemented in a package called `postgres` in `internal/domain/health/repository/postgres/postgres.go`

    func (r *repository) Readiness() error {
    	return r.db.Ping()
    }


### Use Case

This is where all business logic lives. By having repository layer underneath in a separate layer, those functions are reusable in other use case layers.

### Handler

This layer is responsible in handling request from outside world and into the `use case` layer. It does the following:

 1. Parse request into private 'request' struct
 2. Sanitize and validates said struct
 3. Pass into `use case` layer
 4. Process results from `use case` layer and decide how the payload is going to be formatted to
  outside world.
  
Route API are defined in `RegisterHTTPEndPoints` in their respective `register.go` file. 


### Dependency Injection

How does dependency injection happens? It starts with `InitDomains()` method. 

    healthHandler.RegisterHTTPEndPoints(s.router, usecase.NewHealthUseCase(postgres.NewHealthRepository(s.db)))

The repository gets access a pointer to `sql.DB` to perform database operations. This layer also knows nothing of layers above it. `NewBookUseCase` depends on that repository and finally the handler depends on the use case.

## Libraries

Initialization of external libraries are located in `third_party/`

## Testing

### Unit Testing

Unit testing can be run with

    task test
    
Which runs `go test -v ./...`

In repository unit testing, tt makes use of [dockertest](https://github.com/ory/dockertest) from ory that spins up temporary database in a docker to run all repositories.

### End to End Test

Start

    task dockertest

or

    cd docker-test && docker-compose down -v --build && docker-compose up -d
    docker exec -t go8_container_test "/home/appuser/app/e2e"

Stop container

    docker-compose down

# TODO

 - [ ] Complete HTTP integration test
 - [x] better control of json output formatting

# Acknowledgements

 * https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html
 * https://github.com/moemoe89/integration-test-golang
 * https://github.com/george-e-shaw-iv/integration-tests-example
 
# Appendix

## Dev Environment Installation

For Ubuntu:

    sudo apt update && sudo apt install git
    wget https://golang.org/dl/go1.15.6.linux-amd64.tar.gz
    sudo tar -C /usr/local -xzf go1.15.6.linux-amd64.tar.gz
    export PATH=$PATH:/usr/local/go/bin
    echo 'PATH=$PATH:/usr/local/go/bin' >> ~/.profile

    curl -s https://get.docker.com | sudo bash

    sudo apt remove docker docker-engine docker.io containerd runc
    sudo apt update
    sudo apt install -y apt-transport-https ca-certificates curl gnupg-agent software-properties-common
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
    sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
    sudo apt update
    sudo apt install -y docker-ce docker-ce-cli containerd.io
    sudo usermod -aG docker ${USER}
    newgrp docker
    su - ${USER} # or logout and login

    sudo curl -L "https://github.com/docker/compose/releases/download/1.27.4/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
