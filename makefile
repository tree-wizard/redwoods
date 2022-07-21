## Version and Repo info
VERSION ?= 0.0.4
NAME = $(shell basename "`pwd`")
GITBASEURL = github.com/xxyyx
CONTAINER_REPOSITORY=microk8s:32000
PROJECTNAME = $(addprefix ${GITBASEURL}/,${NAME})
IMAGE_NAME=${NAME}:${VERSION}
IMAGE_NAME_LATEST=${NAME}:latest
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

## Golang specific settings
GOOS?=linux
GOARCH?=amd64
GO_BIN := ${GOPATH}/bin
GOVERSION := $(shell go version | awk '{print $$3}')

define HEADER
@echo "###############################################################"
@echo Welcome to the Redwoods fuzzing Suite
@echo "###############################################################"
endef

default: all

help: ## Display this help.
	$(HEADER) 
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

fmt: ## Run go fmt against code.
	go fmt ./...

vet: ## Run go vet against code.
	go vet ./...

test: fmt vet ## Run tests.
	 go test ./...

run: fmt vet ## Run.
	 go run ./cmd/redwoods.go


##@ Core

clean: ## remove previous binaries
	rm -f redwoods

build: clean ## build a version of the app, pass Buildversion, Comit and projectname as build arguments
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build \
		-ldflags "-s -w -X ${PROJECTNAME}/pkg/version.Release=${VERSION} \
		-X ${PROJECTNAME}/pkg/version.Commit=${COMMIT} -X ${PROJECTNAME}/pkg/version.BuildTime=${BUILD_TIME}" \
		-o redwoods ./cmd/redwoods/main.go

install: ## install redwoods into gobin
	cp ./redwoods $(GO_BIN)/redwoods


uninstall: ## uninstall redwoods from gobin
	rm -f $(GO_BIN)/redwoods

##@ Container Deployment
docker-build: ## Build the docker image and tag it with the current version and :latest
		sudo docker build -t ${IMAGE_NAME} -t ${IMAGE_NAME_LATEST} --build-arg tz=${TIMEZONE} . -f ./dockerfiles/Dockerfile

docker-run: docker-build ## Build the docker image and tag it and run it in docker
	sudo docker stop $(IMAGE_NAME) || true && sudo docker rm $(IMAGE_NAME) || true
	sudo docker run --name ${NAME} -v $(shell pwd)/fuzz:/app/redwoods/fuzz/ -v $(shell pwd)/redwoods-cfg.json:/app/redwoods/redwoods-cfg.json --rm \
		$(IMAGE_NAME)  redwoods workflow work 

docker-run-interactive: docker-build  ## run an interactive container
	sudo docker stop $(IMAGE_NAME) || true && sudo docker rm $(IMAGE_NAME) || true
	sudo docker run --name ${NAME} --rm -it  -w '/app/redwoods/'  $(IMAGE_NAME) /bin/bash 

docker-push: ##push your image to the docker hub
	sudo docker tag ${IMAGE_NAME} ${CONTAINER_REPOSITORY}/${IMAGE_NAME}
	sudo docker tag ${IMAGE_NAME_LATEST} ${CONTAINER_REPOSITORY}/${IMAGE_NAME_LATEST}
	sudo docker push  ${CONTAINER_REPOSITORY}/${IMAGE_NAME}
	sudo docker push  ${CONTAINER_REPOSITORY}/${IMAGE_NAME_LATEST}