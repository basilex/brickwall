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

### 3. Playing with an API

The API server opens and maps the port 8081 by default to the localhost

Check the API from the browser and make the requests to:

- http://localhost:8081/api/v1/aux
- http://localhost:8081/api/v1/aux/health
- http://localhost:8081/api/v1/aux/metadata
- http://localhost:8081/api/v1/aux/environment

The rest of endpoints you may check int the source code microservice controllers.

### 4. API Docs

Swagger API docs will be implemented in the nearest future

### 5. NB

_PS:_ Project is in the active development so no concrete instructions or stable structure.

_PS2:_ This readme will be extended, depending of the stage of the development process.

Happy hacking!
