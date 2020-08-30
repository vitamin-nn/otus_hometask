ifndef $(GOPATH)
    GOPATH=$(shell go env GOPATH)
    export GOPATH
endif

build:
	CGO_ENABLED=0 GOARCH=amd64 go build -o ./.bin/calendar

test:
	go test -v -race ./...

install-deps:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(GOPATH)/bin v1.30.0

lint: install-deps
	golangci-lint run ./...

run-compose:
	docker-compose -f ./deployments/docker-compose.yml up -d

down-compose:
	docker-compose -f ./deployments/docker-compose.yml down

run-server:
	source ./configs/local.env && source ./configs/local.server.env && go run . server

run-sc:
	source ./configs/local.env && source ./configs/local.scheduler.env && go run . scheduler

run-sender:
	source ./configs/local.env && source ./configs/local.sender.env && go run . sender

migrate:
	goose -dir migrations postgres "user=vdudov password=123qwe dbname=calendar sslmode=disable" up

generate:
	protoc -Iapi -I/usr/local/include -I$(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.14.6/third_party/googleapis  --go_out=plugins=grpc:internal/server/grpc calendar.proto
	protoc -Iapi -I/usr/local/include -I$(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.14.6/third_party/googleapis --grpc-gateway_out=logtostderr=true:internal/server/grpc calendar.proto
.PHONY: build test lint run-compose down-compose run-server run-sc migrate generate
