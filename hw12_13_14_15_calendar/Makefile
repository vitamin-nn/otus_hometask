ifndef $(GOPATH)
    GOPATH=$(shell go env GOPATH)
    export GOPATH
endif

build:
	go build -o ./.bin/calendar ./cmd/calendar/main.go

test:
	go test -v -race ./...

lint:
	golangci-lint run ./...

run:
	go run ./cmd/calendar/main.go

# далее dsn для постгреса должен уехать в переменную окружения и не будет здесь светиться
migrate:
	goose -dir migrations postgres "user=vdudov password=123qwe dbname=calendar sslmode=disable" up

generate:
	protoc -Iapi -I/usr/local/include -I$(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.14.6/third_party/googleapis  --go_out=plugins=grpc:internal/server/grpc calendar.proto
	protoc -Iapi -I/usr/local/include -I$(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.14.6/third_party/googleapis --grpc-gateway_out=logtostderr=true:internal/server/grpc calendar.proto
.PHONY: build
