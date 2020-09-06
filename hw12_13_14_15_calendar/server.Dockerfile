FROM golang:1.14-alpine AS builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/calendar

FROM scratch

COPY --from=builder /go/bin/calendar /go/bin/calendar
ENTRYPOINT ["/go/bin/calendar"]
CMD ["server"]
EXPOSE 8081 8080
