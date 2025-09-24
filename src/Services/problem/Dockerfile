FROM golang:latest AS base

RUN apt-get update && apt-get install -y --no-install-recommends \
    latexmk \
    texlive-xetex \
    texlive-latex-recommended \
    texlive-latex-extra \
    texlive-fonts-recommended \
    texlive-fonts-extra \
    texlive-lang-all \
    perl \
    wget \
    make \
 && apt-get clean && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Go dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go mod tidy

FROM base AS dev
RUN go install github.com/air-verse/air@latest
RUN mkdir /storage
EXPOSE 3000
CMD ["air", "-c", ".air.toml"]

FROM base AS dev_ssh
WORKDIR /app
RUN CGO_ENABLED=0 GOMAXPROCS=1 go build -p 1 -o /go/bin/main ./cmd/main.go

EXPOSE 3000

CMD ["/go/bin/main"]

FROM base AS builder
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/main.go

RUN --mount=type=cache,target=/root/.cache/go-build \
    GOMAXPROCS=1 \
    go build -o main ./cmd/main.go

FROM builder AS prod
WORKDIR /root/

COPY --from=builder /app/main .
RUN mkdir /storage
EXPOSE 3000

ENTRYPOINT ["./main"]

