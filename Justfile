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
    #!/usr/bin/env sh
    GO_FILES=$(find ./back -type f -name '*.go' ! -path './vendor/*'); 
    if [ -n "$GO_FILES" ]; then 
        go run golang.org/x/tools/cmd/goimports@latest -w $GO_FILES; 
        golines -m 120 -w $GO_FILES;
    fi

# Check backend Go formatting without modifying files
back-format-check:
    #!/usr/bin/env sh
    GO_FILES=$(find ./back -type f -name '*.go' ! -path './vendor/*'); 
    if [ -n "$GO_FILES" ]; then 
        GOIMPORTS_OUT=$(go run golang.org/x/tools/cmd/goimports@latest -l $GO_FILES); 
        if [ -n "$GOIMPORTS_OUT" ]; then
            printf "%s\n" "$GOIMPORTS_OUT";
            printf "goimports found unformatted files.\n";
            exit 1;
        fi;
        GOLINES_OUT=$(golines -m 120 -l $GO_FILES);
        if [ -n "$GOLINES_OUT" ]; then
            printf "%s\n" "$GOLINES_OUT";
            printf "golines found unformatted files.\n";
            exit 1;
        fi;
    fi

# Lint backend Go code
back-lint:
    cd back && golangci-lint run ./...

# Run backend tests
back-test:
    cd back && go test ./... | grep -v \?

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
