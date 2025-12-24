# Isolate
FROM debian:bookworm as sandbox
RUN apt-get update && apt-get install -y wget gcc g++ git build-essential libcap-dev
WORKDIR /app
RUN git clone --depth 1 --branch v1.10.1 https://github.com/ioi/isolate.git
WORKDIR /app/isolate
RUN make isolate

# build starts here
FROM golang:latest AS base

# RUN apt-get update && apt-get install -y openjdk-17-jdk mono-mcs fpc cython3 golang pandoc gccgo pypy3 python3-dev g++ gcc
# Pull isolate binary from sandbox to base
RUN apt-get update && apt-get install -y g++ gcc python3 pypy3 default-jre
COPY --from=sandbox /app/isolate /app/isolate
WORKDIR /app/isolate
RUN make install

# Install go dependencies
WORKDIR /app

COPY go.mod go.sum 
.
RUN go mod download

COPY . .
RUN go mod tidy

# Dev build
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

EXPOSE 8000

ENTRYPOINT ["./main"]
