FROM golang:1.23-alpine as builder

WORKDIR /app

COPY ../go.mod ./

RUN go mod download

COPY .. .

RUN go build -o main ./cmd/app/server.go

# stage2 :running
FROM debian:bullseye-slim


WORKDIR /root/

COPY --from=builder /app/main .


COPY ../config.yml .

EXPOSE 8080

CMD ["./main"]
