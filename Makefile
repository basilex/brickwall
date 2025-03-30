#
# Makefile
# Brickwall Platform Service
# Copyright (C) 2025 by Brickwall Inc. All Rights Reserved
# --------------------------------------------------------
#
# Base environment
#
svc := $(or $(BSP_SVC),bsp)
env := $(or $(BSP_ENV),dev)
ver := $(or $(BSP_VER),0.1.0)
sys := $(or $(BSP_SYS),brickwall)
cert := $(or $(BSP_CERT),./resource/cert)

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
# Check the dependend bins
#
bins = go openssl

-include internal/storage/Makefile.inc

checkfor := $(foreach exec,$(bins), \
	$(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH)))
#
# Main entry point
#
all:
	@echo '*** Help will be implemented later'
	@exit 0

.PHONY: all
#
# Swagger section
#
docs:
	swag init

.PHONY: docs
#
# App section
#
app-build:
	@go build -a -ldflags="$(ldflags)" -o $(svc) main.go
app-tidy:
	@go mod tidy

.PHONY: app-build app-tidy
#
# Compose section
#
app-up:
	@docker compose -f compose-local.yml up --build
app-down:
	@docker compose -f compose-local.yml down  --remove-orphans
app-clean:
	@docker rm -v $(shell docker ps --filter status=exited -q)
	@docker rmi $(img)
app-cert:
	# <dev> mode using self signed certificate
	# <prod> mode using let's encrypt for the real domain
	@openssl req -x509 -newkey rsa:4096 \
		-keyout $(cert)/$(svc).key -out $(cert)/$(svc).crt -days 365 -nodes
app-prune:
	@docker system prune -af

.PHONY: app-up app-down app-clean app-cert app-prune
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

.PHONY: dbs-gen dbs-up dbs-up1 dbs-down dbs-down1 dbs-drop dbs-version
#
# eof
#
