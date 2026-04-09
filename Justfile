set shell := ["bash", "-cu"]

default:
    @just --list

# Install frontend dependencies
front-install:
    corepack pnpm install --dir front

# Run backend API
back-dev:
    cd back && air -c .air.toml

# Format backend Go code
back-format:
    cd back && go fmt ./...

# Check backend Go formatting without modifying files
back-format-check:
    @unformatted=$(cd back && gofmt -l .); if [ -n "$unformatted" ]; then printf 'Unformatted Go files:\n%s\n' "$unformatted"; exit 1; fi

# Lint backend Go code
back-lint:
    cd back && golangci-lint run ./...

# Run backend tests
back-test:
    cd back && go test ./...

# Run backend quality checks
back-check:
    just back-format-check && just back-lint && just back-test

# Run frontend app
front-dev:
    corepack pnpm --dir front dev

# Lint frontend code (non-mutating)
front-lint:
    corepack pnpm --dir front lint

# Format frontend code (mutating)
front-format:
    corepack pnpm --dir front format

# Run frontend tests
front-test:
    corepack pnpm --dir front test

# Run frontend tests in watch mode
front-test-watch:
    corepack pnpm --dir front test:watch

# Run frontend quality checks (canonical non-mutating FE validation entrypoint)
front-check:
    corepack pnpm --dir front format:check && corepack pnpm --dir front lint && corepack pnpm --dir front test

# Run frontend + backend together
dev:
    trap 'kill 0' EXIT; (cd back && air -c .air.toml) & (corepack pnpm --dir front dev) & wait
