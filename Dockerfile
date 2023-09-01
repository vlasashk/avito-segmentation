FROM golang:1.21-alpine as builder
WORKDIR /router
COPY go.mod go.sum ./
RUN go mod download
COPY . /router
RUN go build -o app ./cmd/app/main.go

FROM alpine
WORKDIR /router
COPY --from=builder /router/app .
COPY /pkg/init_sql/*.sql .
COPY /config/config.yaml .

EXPOSE 8090
ENTRYPOINT ["/router/app"]