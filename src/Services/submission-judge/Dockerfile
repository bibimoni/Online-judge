# Base stage - build isolate and Go app in one stage
FROM golang:1.25.5-trixie AS base

# Install all dependencies (build tools + runtime)
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    wget \
    gcc \
    g++ \
    make \
    build-essential \
    libcap-dev \
    libcap2 \
    libcap2-bin \
    libsystemd-dev \
    libsystemd0 \
    pkg-config \
    python3 \
    pypy3 \
    default-jre \
    && rm -rf /var/lib/apt/lists/*

# Download, build, and install isolate v2.2.1
WORKDIR /tmp/isolate-build
RUN wget https://github.com/ioi/isolate/archive/refs/tags/v2.2.1.tar.gz && \
    tar -xzf v2.2.1.tar.gz && \
    cd isolate-2.2.1 && \
    make isolate && \
    make install && \
    cd / && \
    rm -rf /tmp/isolate-build

# Install systemd (for systemd-tmpfiles utility)
RUN apt-get update && apt-get install -y systemd && rm -rf /var/lib/apt/lists/*

# Create systemd-tmpfiles configuration for isolate
RUN mkdir -p /usr/lib/tmpfiles.d && \
    echo 'd /run/isolate 0777 root root -' > /usr/lib/tmpfiles.d/isolate.conf && \
    echo 'd /run/isolate/locks 0777 root root -' >> /usr/lib/tmpfiles.d/isolate.conf

# Create isolate configuration file
RUN mkdir -p /usr/local/etc && \
    echo "box_root = /var/local/lib/isolate" > /usr/local/etc/isolate && \
    echo "lock_root = /run/isolate/locks" >> /usr/local/etc/isolate && \
    echo "cg_root = /sys/fs/cgroup" >> /usr/local/etc/isolate && \
    echo "first_uid = 60000" >> /usr/local/etc/isolate && \
    echo "first_gid = 60000" >> /usr/local/etc/isolate && \
    echo "num_boxes = 100" >> /usr/local/etc/isolate

# Set proper capabilities on isolate binary and create directories
RUN setcap cap_sys_admin=eip /usr/local/bin/isolate && \
    chmod u+s /usr/local/bin/isolate && \
    mkdir -p /var/local/lib/isolate && \
    chmod 777 /var/local/lib/isolate

# Run systemd-tmpfiles on startup to create /run directories
RUN echo '#!/bin/bash' > /entrypoint.sh && \
    echo 'systemd-tmpfiles --create /usr/lib/tmpfiles.d/isolate.conf 2>/dev/null || true' >> /entrypoint.sh && \
    echo 'exec "$@"' >> /entrypoint.sh && \
    chmod +x /entrypoint.sh

# Set up Go workspace
WORKDIR /app

# Copy go mod files first (for better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code (this layer changes most frequently)
COPY . .

# Dev build
FROM base AS dev
RUN go install github.com/air-verse/air@latest
EXPOSE 8000
ENTRYPOINT ["/entrypoint.sh"]
CMD ["air", "-c", ".air.toml"]

# Prod build 
FROM base AS builder
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./src/cmd/main.go

# Prod runtime - minimal image
FROM alpine:latest AS prod
RUN apk add --no-cache libcap libstdc++
WORKDIR /root/
COPY --from=builder /app/main .
RUN chmod +x main
EXPOSE 8000
ENTRYPOINT ["./main"]
