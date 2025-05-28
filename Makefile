# Makefile for Chaos Engineering as a Service

# Variables
BINARY_NAME_CONTROLLER=controller
BINARY_NAME_API=api-server
DOCKER_REPO=chaos-engineering
DOCKER_TAG=latest
GO_BUILD_FLAGS=-v

# Directories
BIN_DIR=bin
DASHBOARD_DIR=dashboard

# Go build targets
.PHONY: build
build: build-controller build-api

.PHONY: build-controller
build-controller:
	mkdir -p $(BIN_DIR)
	go build $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BINARY_NAME_CONTROLLER) ./cmd/controller

.PHONY: build-api
build-api:
	mkdir -p $(BIN_DIR)
	go build $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BINARY_NAME_API) ./api

# Docker build targets
.PHONY: docker-build
docker-build: docker-build-controller docker-build-api

.PHONY: docker-build-controller
docker-build-controller:
	docker build -t $(DOCKER_REPO)/$(BINARY_NAME_CONTROLLER):$(DOCKER_TAG) -f Dockerfile .

.PHONY: docker-build-api
docker-build-api:
	docker build -t $(DOCKER_REPO)/$(BINARY_NAME_API):$(DOCKER_TAG) -f api/Dockerfile .

# Dashboard targets
.PHONY: dashboard-install
dashboard-install:
	cd $(DASHBOARD_DIR) && npm install

.PHONY: dashboard-build
dashboard-build:
	cd $(DASHBOARD_DIR) && npm run build

.PHONY: dashboard-dev
dashboard-dev:
	cd $(DASHBOARD_DIR) && npm start

# Kubernetes deployment targets
.PHONY: deploy
deploy: deploy-crds deploy-controller

.PHONY: deploy-crds
deploy-crds:
	kubectl apply -f deploy/kubernetes/crds/

.PHONY: deploy-controller
deploy-controller:
	kubectl apply -f deploy/kubernetes/deployment.yaml

# Code generation targets
.PHONY: generate
generate:
	./hack/update-codegen.sh

# Clean targets
.PHONY: clean
clean:
	rm -rf $(BIN_DIR)
	rm -rf $(DASHBOARD_DIR)/build
	rm -rf $(DASHBOARD_DIR)/node_modules

# Run targets
.PHONY: run-api
run-api: build-api
	./$(BIN_DIR)/$(BINARY_NAME_API)

# All-in-one target for local development
.PHONY: dev
dev: generate build dashboard-build run-api
