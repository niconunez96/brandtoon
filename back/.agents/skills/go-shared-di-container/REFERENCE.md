# Reference — go-shared-di-container

## Current singleton root

```go
package shared

var container *DIContainer

type DIContainer struct {
    // cached dependencies
}

func NewDIContainer() *DIContainer {
    if container == nil {
        container = &DIContainer{}
    }
    return container
}
```

## Lazy getter pattern

```go
func (c *DIContainer) GetPostgres() *sql.DB {
    if c.postgres == nil {
        c.postgres = mustOpenPostgresFromEnv()
    }
    return c.postgres
}

func (c *DIContainer) GetIdentityUserRepo() identityDomain.UserRepository {
    if c.userRepo == nil {
        c.userRepo = repo.NewUserPostgresRepo(c.GetPostgres())
    }
    return c.userRepo
}
```

## Composition in main

```go
func main() {
    c := shared.NewDIContainer()

    apiDeps := http.Dependencies{
        UserRepo: c.GetIdentityUserRepo(),
    }

    http.RegisterRoutes(api, apiDeps)
}
```

## Design notes

- Keep `DIContainer` in shared context to avoid duplicated initialization logic.
- Use getters to centralize lifecycle and reuse for infra adapters.
- Return domain ports from getters when wiring use cases/handlers.
