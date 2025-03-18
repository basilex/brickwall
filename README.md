# Brickwall Platform

Platform of the RESTful API microservices platform.

## Installation

0. External servies

Install postgres, redis, nats locally or obtain the appropriate connections.

1. Database initialization

Go to internal/storage directory and check the Makefile regarding the migrations (gen, up, down, drop...).
Setup the database, installed somewhere (NOT IN DOCKER) using by the concrete command in the Makefile.

2. Microservices platform building and executing

In the root of the project check the corresponding command in the Makefile.

You've to 'make up' or 'make down' to start or stop the docker compose microservice containers.

PS: Project is in the active development so no concrete instructions or stable structure.
PS2: I'm promise - this readme will be extended, depending of the stage of the development process.

Good luck!

## API Docs

Swagger API docs will be implemented in the nearest future

