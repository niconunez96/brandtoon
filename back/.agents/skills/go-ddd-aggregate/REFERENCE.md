# Reference - Go DDD Aggregate

## Directory Contract

```text
bounded_contexts/{context}/{aggregate}/
  domain/
    aggregate_root.go
    entities.go
    value_objects.go
    ports.go
```

## Aggregate Rules

- Export aggregate root type and constructors.
- Keep child entities private to package when possible.
- Expose behavior methods, not data mutation methods.

## Value Object Pattern

```go
package userdomain

type UserStatus string

const (
    UserStatusPending  UserStatus = "pending"
    UserStatusActive   UserStatus = "active"
    UserStatusBlocked  UserStatus = "blocked"
)
```

Use typed constants in commands, policies, and validations.

## Port (Domain Interface) Pattern

```go
package userdomain

import "context"

type EmailSender interface {
    SendResetPassword(ctx context.Context, to string, token string) error
}

type UserRepository interface {
    Save(ctx context.Context, user *UserAggregate) error
    FindByEmail(ctx context.Context, email string) (*UserAggregate, error)
}
```

Keep signatures infrastructure-agnostic. No SQL types, HTTP objects, or driver-specific contracts.

## Provider-Agnostic Naming (Domain)

Use business language that remains valid if provider changes.

| Avoid (provider-coupled) | Prefer (domain language) |
| --- | --- |
| `GoogleIdentityProvider` | `IdentityProvider` |
| `googleSubject` | `ExternalSubject` or `ProviderSubject` |
| `GoogleUserProfile` | `Identity` or `ExternalIdentity` |

Rule of thumb: if swapping Google with GitHub forces a rename in `domain/`, the naming is wrong.

## Invariant Pattern

```go
func (u *UserAggregate) Activate() error {
    if u.status == UserStatusBlocked {
        return ErrBlockedUserCannotActivate
    }
    u.status = UserStatusActive
    return nil
}
```

State transitions happen through explicit methods only.

## Review Questions

- Can any external package mutate child entities directly?
- Are invariants enforced in aggregate methods?
- Are value objects typed, not free-form strings?
- Do domain interfaces avoid infrastructure dependencies?
