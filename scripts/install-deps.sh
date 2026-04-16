#!/bin/bash
set -e
export PATH="/usr/local/go/bin:/root/go/bin:$PATH"

echo "Installing auth-v2 dependencies..."
(cd /code/auth-v2 && npm install)

echo "Installing gateway dependencies..."
(cd /code/gateway && go mod download)

echo "Installing problem dependencies..."
(cd /code/problem && go mod download)

echo "Installing contest dependencies..."
(cd /code/contest && go mod download)

echo "Installing submission-judge dependencies..."
(cd /code/submission-judge && go mod download)

echo "All dependencies installed."
