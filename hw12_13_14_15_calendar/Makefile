ifndef $(GOPATH)
    GOPATH=$(shell go env GOPATH)
    export GOPATH
endif

up: build-docker-server build-docker-migrate build-docker-scheduler build-docker-sender
	docker-compose -f ./deployments/docker-compose.yml up -d

down:
	docker-compose -f ./deployments/docker-compose.yml down

build:
	CGO_ENABLED=0 GOARCH=amd64 go build -o ./.bin/calendar

build-docker-server:
	docker build -f ./server.Dockerfile -t calendar/server ./

build-docker-scheduler:
	docker build -f ./scheduler.Dockerfile -t calendar/scheduler ./

build-docker-sender:
	docker build -f ./sender.Dockerfile -t calendar/sender ./

build-docker-migrate:
	docker build -t calendar/migrate ./migrate

build-docker-int-tests:
	docker build -f ./integration-test.Dockerfile -t calendar/integration-tests ./

bdd: build-docker-server build-docker-migrate build-docker-scheduler build-docker-sender build-docker-int-tests
	./scripts/run_int_tests.sh

test:
	go test -v -race -gcflags=-l -timeout=30s ./...

bdd-local:
	source ./configs/local.env && source ./configs/local.server.env && cd ./tests/integration/apihttp && godog

install-deps:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(GOPATH)/bin v1.30.0

lint: install-deps
	golangci-lint run ./...

run-server:
	source ./configs/local.env && source ./configs/local.server.env && go run . server

run-sc:
	source ./configs/local.env && source ./configs/local.scheduler.env && go run . scheduler

run-sender:
	source ./configs/local.env && source ./configs/local.sender.env && go run . sender

migrate:
	goose -dir migrate/migrations postgres "user=calendar password=passwd dbname=calendar sslmode=disable" up

generate:
	protoc -Iapi -I/usr/local/include -I$(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.14.6/third_party/googleapis -I$(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.14.6  --go_out=plugins=grpc:internal/server/grpc calendar.proto
	protoc -Iapi -I/usr/local/include -I$(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.14.6/third_party/googleapis -I$(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.14.6 --grpc-gateway_out=logtostderr=true:internal/server/grpc --swagger_out=allow_merge=true,merge_file_name=api:api calendar.proto
	swagger generate client -f ./api/api.swagger.json -t ./tests/httpclient

.PHONY: build test lint run-compose down-compose run-server run-sc migrate generate bdd
