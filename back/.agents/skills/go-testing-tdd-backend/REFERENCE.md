# Reference — go-testing-tdd-backend

## Minimal TDD workflow per change

1. Write/update failing test in target `test/` directory.
2. Implement minimal production code.
3. Re-run tests and make them pass.
4. Refactor only with green tests.

---

## Test locations

```text
bounded_contexts/{context}/{aggregate}/useCases/
  {use_case}.go
  test/
    {use_case}_test.go

bounded_contexts/{context}/{aggregate}/infra/http/
  routes.go
  {handler}.go
  test/
    {handler}_test.go
```

---

## Mock locations (reusable)

```text
bounded_contexts/shared/domain/mocks/
bounded_contexts/{context}/{aggregate}/domain/mocks/
```

Use aggregate-local mocks for aggregate interfaces. Use shared mocks only for truly shared interfaces.

---

## Example reusable mock factory

```go
package mocks

type MockUserRepository struct {
    SaveErr error
}

type MockUserRepositoryOption func(*MockUserRepository)

func NewMockUserRepository(opts ...MockUserRepositoryOption) *MockUserRepository {
    m := &MockUserRepository{}
    for _, opt := range opts {
        opt(m)
    }
    return m
}

func WithSaveErr(err error) MockUserRepositoryOption {
    return func(m *MockUserRepository) {
        m.SaveErr = err
    }
}
```

---

## Use case test intent

- Validate command/query outcomes.
- Validate side effects at observable boundary.
- Assert errors precisely.

## Handler test intent

- Validate HTTP status/body/headers.
- Validate request-to-usecase mapping.
- Avoid embedding domain decision logic in handler tests.

---

## Huma-specific handler test pattern

Reference: https://huma.rocks/tutorial/writing-tests/

Prefer Huma test server utilities for endpoint-level tests:

```go
func TestCreateThing(t *testing.T) {
    _, api := humatest.New(t)
    RegisterRoutes(api, deps)

    resp := api.Post("/things", map[string]any{
        "name": "demo",
    })

    if resp.Code != 201 {
        t.Fatalf("unexpected status: %d", resp.Code)
    }
}
```

If route behavior depends on request context, use context-aware methods like `GetCtx`/`PostCtx`.
