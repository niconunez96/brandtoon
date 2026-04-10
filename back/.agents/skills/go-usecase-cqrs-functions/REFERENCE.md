# Reference - Go Use Cases with CQRS Functions

## Directory Contract

```text
bounded_contexts/{context}/{aggregate}/useCases/
  authenticate_user.go
  reset_password.go
  get_user_profile.go
```

One business action per file.

## Command Template

```go
package useCases

import (
    "context"
    "brandtoon/back/bounded_contexts/identity/user/domain"
)

type ResetPasswordCommand struct {
    UserID string
    Email  string
}

func ResetPassword(
    ctx context.Context,
    cmd ResetPasswordCommand,
    users domain.UserRepository,
    sender domain.EmailSender,
) error {
    user, err := users.FindByEmail(ctx, cmd.Email)
    if err != nil {
        return err
    }

    token, err := user.GenerateResetToken()
    if err != nil {
        return err
    }

    if err := users.Save(ctx, user); err != nil {
        return err
    }

    return sender.SendResetPassword(ctx, cmd.Email, token)
}
```

Command must perform state change and/or side effect.

## Query Template

```go
package useCases

import (
    "context"
    "brandtoon/back/bounded_contexts/identity/user/domain"
)

type GetUserProfileQuery struct {
    Email string
}

type UserProfileView struct {
    Email  string
    Status string
}

func GetUserProfile(
    ctx context.Context,
    q GetUserProfileQuery,
    users domain.UserRepository,
) (UserProfileView, error) {
    user, err := users.FindByEmail(ctx, q.Email)
    if err != nil {
        return UserProfileView{}, err
    }

    return UserProfileView{
        Email:  user.Email(),
        Status: string(user.Status()),
    }, nil
}
```

Query must avoid writes/side effects.

## Provider-Agnostic Naming (Use Cases)

Use case names should describe business intent and remain stable across provider swaps.

| Avoid (provider-coupled) | Prefer (business/domain) |
| --- | --- |
| `AuthenticateGoogleCallback` | `AuthenticateCallback` |
| `GetGoogleAuthURLQuery` | `GetAuthURLQuery` |
| `GoogleAuthURLResult` | `AuthURLResult` |

Quick heuristic:
- If you can replace provider implementation in `infra/` without touching `useCases/`, naming is correct.
- If names need to change when provider changes, naming is too coupled.

## Review Questions

- Does each file contain exactly one use case?
- Is function signature explicit about dependencies?
- Does command/query behavior match its name?
- Is transport logic absent from use case layer?
