#
# Makefile
# Brickwall SaaS Platform Service
#
# Base environment
#
svc := $(or $(BSP_SVC),bsp)
env := $(or $(BSP_ENV),dev)
ver := $(or $(BSP_VER),0.1.0)
sys := $(or $(BSP_SYS),brickwall)

img := $(or $(BSP_IMG),$(sys)/$(svc):$(ver))
#
# Variables to be embedded
#
version := $(ver)
staging := $(env)
githash := $(shell git rev-parse --short=8 HEAD)
gobuild := $(shell go version | sed -e "s/go version //g;s/ /-/g")
compile := $(shell date "+%FT%T.%N%:z")
#
# Go build: ldflags linker parameters
#
ldflags += -X main.Version=$(version)
ldflags += -X main.Staging=$(staging)
ldflags += -X main.Githash=$(githash)
ldflags += -X main.Gobuild=$(gobuild)
ldflags += -X main.Compile=$(compile)

all: build

build:
	@go build -a -ldflags="$(ldflags)" -o $(svc) main.go

up:
	@docker compose up --build --force-recreate

down:
	@docker compose down  --remove-orphans

clean:
	@docker rm -v $(shell docker ps --filter status=exited -q)
	@docker rmi $(img)

prune:
	@docker system prune -f

tidy:
	@go mod tidy

.PHONY: all build up down clean prune tidy
#
# eof
#
