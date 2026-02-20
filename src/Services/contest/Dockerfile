FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go mod tidy

RUN go install github.com/air-verse/air@latest
EXPOSE 8001
# RUN air init
# CMD ["go", "run", "cmd/main.go"]
CMD ["air", "-c", ".air.toml"]
