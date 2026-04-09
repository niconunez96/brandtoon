set shell := ["bash", "-cu"]

default:
    @just --list

# Install frontend dependencies
front-install:
    pnpm install --dir front

# Run backend API
back-dev:
    cd back && air -c .air.toml

# Format backend Go code
back-format:
    cd back && go fmt ./...

# Lint backend Go code
back-lint:
    cd back && golangci-lint run ./...

# Run backend tests
back-test:
    cd back && go test ./...

# Run backend quality checks
back-check:
    just back-format && just back-lint && just back-test

# Run frontend app
front-dev:
    pnpm --dir front dev

# Run frontend + backend together
dev:
    trap 'kill 0' EXIT; (cd back && air -c .air.toml) & (pnpm --dir front dev) & wait
