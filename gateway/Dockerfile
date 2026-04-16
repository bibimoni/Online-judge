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

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./src/cmd/main.go

# prod runtime
FROM alpine:latest AS prod
WORKDIR /root/

# COPY .env .
COPY --from=builder /app/main .

RUN chmod +x main

EXPOSE 80

ENTRYPOINT ["./main"]

