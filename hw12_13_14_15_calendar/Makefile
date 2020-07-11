ifndef $(GOPATH)
    GOPATH=$(shell go env GOPATH)
    export GOPATH
endif

build:
	go build -o ./.bin/calendar ./cmd/calendar/main.go

test:
	go test -race ./...

install-deps:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh |  sh -s v1.27.

lint: install-deps
	golangci-lint run ./...

run:
	go run ./cmd/calendar/main.go

# далее dsn для постгреса должен уехать в переменную окружения и не будет здесь светиться
migrate:
	goose -dir migrations postgres "user=vdudov password=123qwe dbname=calendar sslmode=disable" up

.PHONY: build
