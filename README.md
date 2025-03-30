# Brickwall Platform

Platform of the RESTful API microservices with NATS centralized messaging.

## Installation

### 0. Create environment

Copy _.env.example_ file to your own _.env_ file.

```
$ cp .env.example .env
```

and make the corresponding modifications according to your environment.

### 1. Database initialization

Check the Makefile regarding the migrations (dbs-gen, dbs-up, dbs-down, dbs-drop...).
Setup the database, using by the concrete _dbs-\*_ command in the root Makefile.

#### Generate database layer

```
$ make dbs-gen
```

#### Create database schema with default values

```
$ make dbs-up
```

#### Down/Drop database schema

```
$ make dbs-down
```

or

```
$ make dbs-drop
```

### 2. Microservices platform building and executing

In the root of the project check the corresponding command in the Makefile.

You've to _make app-up_ or _make app-down_ to start or stop the docker compose microservice containers.

#### Start docker composer microservice containers

```
$ make app-up
```

#### Stop composer microservice containers

```
$ make app-down
```

### 3. Review Makefile commands

You make cal _make_ without flags to obtain extended help for all the commands available for project managing.

```
$ make

*** Brickwall Makefile sections
    ---------------------------
>>> dbs management section
  - dbs-gen     : Generate sqlc db layer
  - dbs-up      : Install db schema and default data
  - dbs-up1     : Migrate one level of the db schema
  - dbs-down    : Uninstall db schema (all the data purged)
  - dbs-down1   : Migrate down one level of the db schema
  - dbs-drop    : Drop entire db schema (all the data purged)
  - dbs-version : Show the db migration version

>>> app management section
  - app-tidy    : Ensure that all imports are satisfied
  - app-build   : Build the application inside the linux container
  - app-up      : Run all the containers from the docker composer yml
  - app-down    : Shut down all the docker compose containers
  - app-clean   : Remove all the docker exited containers
  - app-cert    : Generate app TLS/SSL certificates
  - app-prune   : Prune all in the local docker env

>>> docker swarm management section
  - stack-load  : Load configs/secrets to docker registry
  - stack-clean : Clean all the configs/secrets from docker registry
  - stack-deploy: Deploy containers to the docker swarm stack
  - stack-remove: Completely remove docker stack (with volumes)
```

### 4. Playing with an API

The API server opens and maps the port 8081 by default to the localhost

Check the API from the browser and make the requests to:

- http://localhost:8081/api/v1/aux
- http://localhost:8081/api/v1/aux/health
- http://localhost:8081/api/v1/aux/metadata
- http://localhost:8081/api/v1/aux/environment

The rest of endpoints you may check int the source code microservice controllers.

### 5. API Docs

Swagger API docs will be implemented in the nearest future

### 6. NB

_PS:_ Project is in the active development so no concrete instructions or stable structure.

_PS2:_ This readme will be extended, depending on the stage of the development process.

_PS3:_ All questions you may send to _alexander.vasilenko@gmail.com_ email.

Happy hacking!
