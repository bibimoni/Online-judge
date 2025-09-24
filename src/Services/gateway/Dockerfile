FROM golang:latest AS base

WORKDIR /app

COPY go.mod go.sum .
RUN go mod download

COPY . .
RUN go mod tidy

# dev build
FROM base AS dev
# COPY .env .
RUN go install github.com/air-verse/air@latest
EXPOSE 8000
CMD ["air", "-c", ".air.toml"]

# prod build 
FROM base AS builder

# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./src/cmd/main.go
RUN --mount=type=cache,target=/root/.cache/go-build \
    GOMAXPROCS=1 \
    go build -o main ./src/cmd/main.go

# prod runtime
# FROM alpine:latest AS prod
# WORKDIR /root/
#
# # COPY .env .
# COPY --from=builder /app/main .
#
# RUN chmod +x main
#
# EXPOSE 81
#
# ENTRYPOINT ["./main"]
#

FROM base AS dev_ssh
WORKDIR /app
RUN CGO_ENABLED=0 GOMAXPROCS=1 go build -p 1 -o /go/bin/main ./src/cmd/main.go
EXPOSE 81
CMD ["/go/bin/main"]

FROM builder AS prod 

WORKDIR /root/
COPY --from=builder /app/main .
RUN chmod +x main

EXPOSE 81
ENTRYPOINT ["./main"]

