# Brickwall Platform

Platform of the RESTful API microservices platform.

## Installation

1. Database initialization

Go to internal/storage directory and check the Makefile regarding the migrations (gen, up, down, drop...).
Setup the database, using by the concrete dbs-* command in the root Makefile.

2. Microservices platform building and executing

In the root of the project check the corresponding command in the Makefile.

You've to 'make api-up' or 'make api-down' to start or stop the docker compose microservice containers.

PS: Project is in the active development so no concrete instructions or stable structure.
PS2: I'm promise - this readme will be extended, depending of the stage of the development process.

Good luck!

## API Docs

Swagger API docs will be implemented in the nearest future

