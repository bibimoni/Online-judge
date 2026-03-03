FROM golang:1.25

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
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
    && apt-get clean && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go mod tidy

RUN ["mkdir", "/storage"]
RUN go install github.com/air-verse/air@latest
EXPOSE 3000
# RUN air init
# CMD ["go", "run", "cmd/main.go"]
CMD ["air", "-c", ".air.toml"]