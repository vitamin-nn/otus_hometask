build:
	go build -o ./.bin/calendar ./cmd/calendar/main.go

test:
	go test -race ./...

lint:
	golangci-lint run ./...

run:
	go run ./cmd/calendar/main.go

# далее dsn для постгреса должен уехать в переменную окружения и не будет здесь светиться
migrate:
	goose -dir migrations postgres "user=vdudov password=123qwe dbname=calendar sslmode=disable" up

.PHONY: build
