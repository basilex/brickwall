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
	@docker compose -f compose.yml up --build
app-down:
	@docker compose -f compose.yml down  --remove-orphans
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
# K8s all sections
#
apply-all:
	kubectl apply -f ./kubernetes/postgres-pvc.yml
	kubectl apply -f ./kubernetes/postgres-config.yml
	kubectl apply -f ./kubernetes/postgres-secrets.yml
	kubectl apply -f ./kubernetes/deployment-postgres.yml
	kubectl apply -f ./kubernetes/deployment-nats.yml
	kubectl apply -f ./kubernetes/deployment-redis.yml
	kubectl apply -f ./kubernetes/deployment-maildev.yml
	kubectl apply -f ./kubernetes/deployment-api.yml
#
# K8s postgres section
#
apply-postgres:
	kubectl apply -f ./kubernetes/postgres-pvc.yml
	kubectl apply -f ./kubernetes/postgres-config.yml
	kubectl apply -f ./kubernetes/postgres-secrets.yml
	kubectl apply -f ./kubernetes/deployment-postgres.yml
redeploy-postgres:
	kubectl rollout restart deployment postgres
delete-postgres:
	kubectl delete deployment postgres
#
# K8s redis section
#
apply-redis:
	kubectl apply -f ./kubernetes/redis-pvc.yml
	kubectl apply -f ./kubernetes/redis-config.yml
	kubectl apply -f ./kubernetes/deployment-redis.yml
redeploy-redis:
	kubectl rollout restart deployment redis
delete-redis:
	kubectl delete deployment redis
#
# .PHONY section
#
.PHONY: all \
	app-docs
	app-up app-down app-clean app-prune app-tidy \
	dbs-gen dbs-up dbs-up1 dbs-down dbs-down1 dbs-drop dbs-version
#
# eof
#
