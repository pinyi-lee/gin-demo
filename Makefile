GO_PATH := $(shell go env GOPATH)
FILE_PATH := $(shell pwd)
SOURCE_PATH := $(FILE_PATH:$(GO_PATH)%=%)

build_gin:
	@rm -f build/Dockerfile
	@export SOURCE_PATH=$(SOURCE_PATH); envsubst '$$SOURCE_PATH' < build/Dockerfile.temp > build/Dockerfile;
	@docker build -t gin_service -f build/Dockerfile .

start:
	@docker-compose -f build/docker-compose.yml up -d

stop:
	@docker-compose -f build/docker-compose.yml stop

check:
	@. app/env.sh 
	@docker-compose -f test/docker-compose.yml up -d
	@go test ./test 
	@sleep 5
	@docker-compose -f test/docker-compose.yml down
