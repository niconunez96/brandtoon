# Reference - Go Infra HTTP and Repository Adapters

Examples in this file complement `back/AGENTS.md`; canonical backend policy lives there.

## Directory Contract

```text
bounded_contexts/{context}/{aggregate}/infra/
  repo/
    user_postgres_repo.go
  http/
    routes.go
    authenticate_user_handler.go
    reset_password_handler.go
```

## Main Wiring Flow

`main.go` composes the app using the shared DI container, creates the Huma API, and injects adapters into aggregate route registration.

```go
func main() {
    c := shared.NewDIContainer()
    userRepo := c.GetIdentityUserRepo()

    router := chi.NewMux()
    api := humachi.New(router, huma.DefaultConfig("Brandtoon API", "1.0.0"))

    identityhttp.RegisterRoutes(api, identityhttp.Dependencies{
        Users: userRepo,
    })
}
```

For the full dependency wiring standard (singleton getters, lifecycle, composition rules), use skill `go-shared-di-container`.

## Routes Pattern

`infra/http/routes.go` registers Huma operations for the aggregate and composes handlers from explicit dependencies.

```go
func RegisterRoutes(api huma.API, deps Dependencies) {
    huma.Register(api, huma.Operation{
        OperationID: "post-authenticate-user",
        Method:      http.MethodPost,
        Path:        "/identity/authenticate",
        Summary:     "Authenticate user",
    }, PostAuthenticateHandler(deps))
}
```

## Handler Pattern (Function-Based)

```go
func PostAuthenticateHandler(deps Dependencies) func(context.Context, *AuthenticateInput) (*AuthenticateOutput, error) {
    return func(ctx context.Context, input *AuthenticateInput) (*AuthenticateOutput, error) {
        cmd := identityusecases.AuthenticateUserCommand{
            Email:    input.Body.Email,
            Password: input.Body.Password,
        }

        out, err := identityusecases.AuthenticateUser(ctx, cmd, deps.Users)
        if err != nil {
            return nil, mapError(err)
        }

        return &AuthenticateOutput{
            Body: out,
        }, nil
    }
}
```

Handlers adapt Huma transport DTOs and call use cases. They do not implement business rules.

## Repository Naming and Contract

- Aggregate infra package naming:
  - `infra/repo/` -> `package {aggregate}repo` (example: `userrepo`)
  - `infra/http/` -> `package {aggregate}http` (example: `userhttp`)
- Implementation naming: `UserPostgresRepo`, `SessionPostgresRepo`, etc.
- Struct must satisfy a domain interface explicitly.

```go
var _ domain.UserRepository = (*UserPostgresRepo)(nil)
```

## Provider-Specific Naming Boundary

Infra is the correct place for concrete provider names.

| Layer | Naming rule | Example |
| --- | --- | --- |
| `infra/` | Provider/vendor naming allowed | `GoogleOAuthClient`, `/auth/google/login` |
| `useCases/` | Business/domain naming only | `AuthenticateCallback`, `GetAuthURLQuery` |
| `domain/` | Business/domain naming only | `IdentityProvider`, `ExternalIdentity` |

Use infra to adapt concrete technologies; keep core language provider-agnostic.

## Review Questions

- Are all third-party imports isolated to infra?
- Are third-party clients initialized in shared `DIContainer` instead of handlers/useCases?
- Are route prefixes centralized in `routes.go`?
- Do handlers call use cases instead of embedding logic?
- Do repo names follow `XXXPostgresRepo` exactly?
- Do `GetX()` container methods implement lazy singleton initialization?
