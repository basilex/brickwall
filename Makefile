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
app-docs:
	swag init
#
# App section
#
app-build:
	@go build -a -ldflags="$(ldflags)" -o $(svc) main.go
app-up:
	@docker compose up --build # --force-recreate
app-down:
	@docker compose down  --remove-orphans
app-clean:
	@docker rm -v $(shell docker ps --filter status=exited -q)
	@docker rmi $(img)
app-prune:
	@docker system prune -f
app-tidy:
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
# Swarm section
#
swarm-init:
	@docker swarm init
swarm-leave:
	@docker swarm leave --force
swarm-setup:
	@chmod +x setenv.sh && ./setenv.sh
swarm-cleanup:
	@docker config ls --format '{{.ID}}' | xargs -r docker config rm
	@docker secret ls --format '{{.ID}}' | xargs -r docker secret rm

swarm-reset: swarm-cleanup swarm-setup

swarm-config-ls:
	@docker config ls
swarm-secret-ls:
	@docker secret ls
swarm-node-ls:
	@docker node ls
#
# .PHONY section
#
.PHONY: all \
	api-docs
	api-up api-down api-clean api-prune api-tidy \
	dbs-gen dbs-up dbs-up1 dbs-down dbs-down1 dbs-drop dbs-version \
	swarm-init swarm-leave swarm-setup swarm-cleanup swarm-reset swarm-config-ls swarm-secret-ls swarm-node-ls
#
# eof
#
