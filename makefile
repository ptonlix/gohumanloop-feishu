# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
PROJECT_NAME=GoHumanLoopWeWork
BINARY_NAME=gohumanloop-wework
BINARY_UNIX=$(BINARY_NAME)_unix
CONFPATH=./conf/app.conf
SHPATH=./run.sh

APP_NAME=${BINARY_NAME}
APP_VERSION= v0.1.1
BUILD_VERSION=$(shell git log -1 --oneline)
BUILD_TIME=$(shell date )
GIT_REVISION=$(shell git rev-parse --short HEAD)
GIT_BRANCH=$(shell git name-rev --name-only HEAD)
GO_VERSION=$(shell go version)
VERSIONINFO = "-s -X 'main.AppName=${APP_NAME}' \
            -X 'main.AppVersion=${APP_VERSION}' \
            -X 'main.BuildVersion=${BUILD_VERSION}' \
            -X 'main.BuildTime=${BUILD_TIME}' \
            -X 'main.GitRevision=${GIT_REVISION}' \
            -X 'main.GitBranch=${GIT_BRANCH}' \
            -X 'main.GoVersion=${GO_VERSION}'" \

define packconfig
	@echo "package config ......."
	mkdir $(PROJECT_NAME)
	mkdir $(PROJECT_NAME)/conf
	mkdir $(PROJECT_NAME)/log
	mkdir $(PROJECT_NAME)/data
	mv $(BINARY_NAME) ./$(PROJECT_NAME)/
	cp $(CONFPATH) $(PROJECT_NAME)/conf/
	cp -r $(SHPATH) $(PROJECT_NAME)/
	tar -zcvf ${BINARY_NAME}-$(APP_VERSION).tar.gz $(PROJECT_NAME)
	rm -rf $(PROJECT_NAME)
endef
all: test build
build:
	$(GOBUILD) -ldflags $(VERSIONINFO) -o $(BINARY_NAME) -v
test:
	$(GOTEST) -v
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v 
	./$(BINARY_NAME)
beego:
	bee run

package: 
	$(GOBUILD) -ldflags $(VERSIONINFO) -o $(BINARY_NAME) -v 
	$(call packconfig)
	@echo "package complete!"

package-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags $(VERSIONINFO) -o $(BINARY_NAME) -v
	$(call packconfig)
	@echo "package complete!"

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
docker-build:
	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v