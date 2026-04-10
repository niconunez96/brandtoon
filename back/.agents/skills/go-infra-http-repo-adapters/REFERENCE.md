# Reference - Go Infra HTTP and Repository Adapters

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

`main.go` composes the app using the shared DI container and injects adapters into route registration.

```go
func main() {
    c := shared.NewDIContainer()
    userRepo := c.GetIdentityUserRepo()

    r := http.NewRouter()
    http.RegisterIdentityRoutes(r, userRepo)
    r.ListenAndServe()
}
```

For the full dependency wiring standard (singleton getters, lifecycle, composition rules), use skill `go-shared-di-container`.

## Routes Pattern

`infra/http/routes.go` groups by prefix and composes handlers.

```go
func RegisterIdentityRoutes(r Router, users domain.UserRepository) {
    grp := r.Group("/identity")
    grp.POST("/authenticate", AuthenticateUserHandler(users))
    grp.POST("/reset-password", ResetPasswordHandler(users))
}
```

## Handler Pattern (Function-Based)

```go
func AuthenticateUserHandler(users domain.UserRepository) HandlerFunc {
    return func(ctx Context) error {
        var req AuthenticateUserRequest
        if err := ctx.Bind(&req); err != nil {
            return ctx.JSON(400, ErrorResponse("invalid payload"))
        }

        cmd := useCases.AuthenticateUserCommand{Email: req.Email, Password: req.Password}
        out, err := useCases.AuthenticateUser(ctx.RequestContext(), cmd, users)
        if err != nil {
            return mapError(ctx, err)
        }

        return ctx.JSON(200, out)
    }
}
```

Handlers adapt transport DTOs and call use cases. They do not implement business rules.

## Repository Naming and Contract

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
