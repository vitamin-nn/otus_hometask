FROM golang:1.14-alpine

#WORKDIR /app
RUN apk add git
RUN apk add gcc
RUN apk add g++
RUN GO111MODULE=on go get github.com/cucumber/godog/cmd/godog@v0.10.0

RUN mkdir -p /opt/go/app/
COPY . /opt/go/app/

WORKDIR /opt/go/app/tests/integration/apihttp

ENTRYPOINT ["godog"]
