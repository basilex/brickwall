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
#
# Main entry point
#
all:
	@echo '*** Help will be implemented later'
	@exit 0
#
# Swagger section
#
api-docs:
	swag init
#
# Service section
#
api-build:
	@go build -a -ldflags="$(ldflags)" -o $(svc) main.go
api-up:
	@docker compose up --build --force-recreate
api-down:
	@docker compose down  --remove-orphans
api-clean:
	@docker rm -v $(shell docker ps --filter status=exited -q)
	@docker rmi $(img)
api-prune:
	@docker system prune -f
api-tidy:
	@go mod tidy
#
# Dbs section
#
dbs-gen:
	@make -C internal/storage gen
dbs-up:
	@make -C internal/storage up
dbs-up1:
	@make -C internal/storage up1
dbs-down:
	@make -C internal/storage down
dbs-down1:
	@make -C internal/storage down1
dbs-drop:
	@make -C internal/storage drop
dbs-version:
	@make -C internal/storage version
#
# PHONY section
#
.PHONY: all \
	api-docs
	api-up api-down api-clean api-prune api-tidy \
	dbs-gen dbs-up dbs-up1 dbs-down dbs-down1 dbs-drop dbs-version
#
# eof
#
