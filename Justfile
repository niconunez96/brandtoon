set shell := ["bash", "-cu"]

default:
    @just --list

# Install frontend dependencies
front-install:
    pnpm install --dir front

# Run backend API
back-dev:
    cd back && go run ./main.go

# Run frontend app
front-dev:
    pnpm --dir front dev

# Run frontend + backend together
dev:
    trap 'kill 0' EXIT; (cd back && go run ./main.go) & (pnpm --dir front dev) & wait
