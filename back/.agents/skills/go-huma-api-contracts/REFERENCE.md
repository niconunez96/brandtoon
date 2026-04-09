# Reference — go-huma-api-contracts

This reference adapts Huma docs to Brandtoon backend conventions.

## Sources

- Service configuration tutorial: https://huma.rocks/tutorial/service-configuration/
- Conditional fields guide: https://huma.rocks/how-to/conditional-fields/
- Writing tests tutorial: https://huma.rocks/tutorial/writing-tests/

---

## Canonical route registration split (Brandtoon style)

Keep operation wiring out of `main` and expose a reusable registration function.

```go
// infra/http/routes.go
package http

import (
    "net/http"

    "github.com/danielgtaylor/huma/v2"
)

type Dependencies struct {
    // Inject use case function deps / services here (interfaces)
}

func RegisterRoutes(api huma.API, deps Dependencies) {
    huma.Register(api, huma.Operation{
        OperationID: "post-authenticate-user",
        Method:      http.MethodPost,
        Path:        "/auth/login",
        Summary:     "Authenticate user",
        Tags:        []string{"Auth"},
    }, PostAuthenticateHandler(deps))
}
```

```go
// infra/http/post_authenticate_handler.go
package http

import (
    "context"
)

type AuthenticateInput struct {
    Body struct {
        Email    string `json:"email" format:"email"`
        Password string `json:"password" minLength:"8"`
    }
}

type AuthenticateOutput struct {
    Body struct {
        AccessToken string `json:"accessToken"`
    }
}

func PostAuthenticateHandler(deps Dependencies) func(context.Context, *AuthenticateInput) (*AuthenticateOutput, error) {
    return func(ctx context.Context, input *AuthenticateInput) (*AuthenticateOutput, error) {
        // map input -> use case cmd/query
        // call use case
        // map result -> output
        out := &AuthenticateOutput{}
        out.Body.AccessToken = "token"
        return out, nil
    }
}
```

---

## Service configuration pattern

Use options struct + CLI/env configuration pattern so runtime settings are explicit and testable.

```go
type Options struct {
    Port int `help:"Port to listen on" short:"p" default:"8888"`
}
```

Compose dependencies via shared container in `main.go`:

```go
func main() {
    router := chi.NewMux()
    api := humachi.New(router, huma.DefaultConfig("Brandtoon API", "1.0.0"))

    deps := http.Dependencies{}
    http.RegisterRoutes(api, deps)
}
```

Dependency wiring details (shared container, singleton getters, composition rules) belong to `go-shared-di-container`.

When both env and CLI args are provided, CLI args take precedence (Huma behavior).

---

## Conditional fields in contracts

When one field requires others, use `dependentRequired` tag in request input:

```go
type Input struct {
    Value      string `json:"value,omitempty" dependentRequired:"dependent1,dependent2"`
    Dependent1 string `json:"dependent1,omitempty"`
    Dependent2 string `json:"dependent2,omitempty"`
}
```

This declares JSON Schema `dependentRequired` for generated docs/validation.

Use this when contract semantics are conditional (e.g., providing `documentType` requires `documentNumber`).

---

## Huma testing setup

Extract route registration into reusable function and use `humatest` in tests:

```go
func addRoutes(api huma.API) {
    // register operations
}

func TestEndpoint(t *testing.T) {
    _, api := humatest.New(t)
    addRoutes(api)
    resp := api.Get("/path")
    if resp.Code != 200 {
        t.Fatalf("unexpected status: %d", resp.Code)
    }
}
```

For context-dependent routes, use `GetCtx` / `PostCtx` variants.

Validation-style test example (422 for invalid payload):

```go
func TestPostAuthenticateValidation(t *testing.T) {
    _, api := humatest.New(t)
    RegisterRoutes(api, Dependencies{})

    // Missing password + invalid email
    resp := api.Post("/auth/login", map[string]any{
        "email": "not-an-email",
    })

    if resp.Code != 422 {
        t.Fatalf("unexpected status: %d", resp.Code)
    }
}
```

---

## Brandtoon mapping notes

- Keep Huma contracts in `infra/http`.
- Map request structs to use case command/query DTOs.
- Keep use cases and domain Huma-agnostic.
