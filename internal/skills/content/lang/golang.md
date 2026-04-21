
# Go Best Practices


## When to Use

- Writing or reviewing any Go file
- Designing Go packages and interfaces
- Handling concurrency with goroutines
- Code review of idiomatic Go patterns

**Don't use when:** The project uses generated code (e.g., protobuf) that overrides standard patterns.

---

## Decision Tree

```text
Function can fail?                     → return (T, error) as last value
Cleaning up a resource?                → defer close/release
Defining shared behavior?              → small interface (1–3 methods)
Sharing data across goroutines?        → channel or sync.Mutex
Cancellable / timeout operation?       → pass context.Context as first arg
Writing a test with multiple cases?    → table-driven test
```

---

## Critical Patterns

### 1. Always Handle Errors Explicitly

```go
// ✅ GOOD — wrap errors with context
user, err := repo.FindByID(ctx, id)
if err != nil {
    return fmt.Errorf("findUser %s: %w", id, err)
}

// ✅ GOOD — sentinel errors for callers to inspect
var ErrNotFound = errors.New("not found")

if errors.Is(err, ErrNotFound) {
    http.NotFound(w, r)
    return
}

// ❌ BAD — ignored error
user, _ := repo.FindByID(ctx, id)
```

### 2. Use `defer` for Cleanup

```go
// ✅ GOOD — guaranteed cleanup even on panic
func readFile(path string) ([]byte, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, fmt.Errorf("open: %w", err)
    }
    defer f.Close()

    return io.ReadAll(f)
}

// ❌ BAD — Close() skipped if an error occurs above it
func readFile(path string) ([]byte, error) {
    f, _ := os.Open(path)
    data, err := io.ReadAll(f)
    f.Close() // not reached on early return
    return data, err
}
```

### 3. Keep Interfaces Small

```go
// ✅ GOOD — single-method interface is easy to satisfy and mock
type UserReader interface {
    FindByID(ctx context.Context, id string) (*User, error)
}

// ✅ GOOD — compose interfaces when needed
type UserRepository interface {
    UserReader
    UserWriter
}

// ❌ BAD — fat interface forces implementors to provide everything
type UserRepository interface {
    FindByID(ctx context.Context, id string) (*User, error)
    FindByEmail(ctx context.Context, email string) (*User, error)
    Create(ctx context.Context, u *User) error
    Update(ctx context.Context, u *User) error
    Delete(ctx context.Context, id string) error
    Count(ctx context.Context) (int, error)
}
```

### 4. Always Pass `context.Context` as First Parameter

```go
// ✅ GOOD — context enables cancellation and timeout propagation
func (s *UserService) GetUser(ctx context.Context, id string) (*User, error) {
    return s.repo.FindByID(ctx, id)
}

// ✅ GOOD — set a timeout at the boundary
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()
user, err := svc.GetUser(ctx, id)

// ❌ BAD — no way to cancel or time out the operation
func (s *UserService) GetUser(id string) (*User, error) {
    return s.repo.FindByID(id)
}
```

### 5. Use Table-Driven Tests

```go
// ✅ GOOD
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
    }{
        {"valid email",  "user@example.com", false},
        {"missing @",    "userexample.com",  true},
        {"empty string", "",                 true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateEmail(tt.email)
            if (err != nil) != tt.wantErr {
                t.Errorf("got err=%v, wantErr=%v", err, tt.wantErr)
            }
        })
    }
}

// ❌ BAD — repeated test functions for each case
func TestValidateEmail_Valid(t *testing.T)   { ... }
func TestValidateEmail_Missing(t *testing.T) { ... }
```

### 6. Use Goroutines Safely

```go
// ✅ GOOD — use errgroup for concurrent fan-out
import "golang.org/x/sync/errgroup"

g, gCtx := errgroup.WithContext(ctx)
g.Go(func() error { return fetchUser(gCtx, id) })
g.Go(func() error { return fetchPosts(gCtx, id) })
if err := g.Wait(); err != nil {
    return err
}

// ✅ GOOD — protect shared state with mutex
type Counter struct {
    mu    sync.Mutex
    value int
}
func (c *Counter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

// ❌ BAD — goroutine leak: no termination condition
go func() {
    for { process() }
}()
```

### 7. Prefer Explicit Fields Over Wide Embedding

```go
// ✅ GOOD — explicit field keeps public API intentional
type Service struct {
    repo   UserRepository
    logger *slog.Logger
}

// ❌ BAD — embedding leaks all repo methods into Service's API
type UserService struct {
    *PostgresRepo // exposes Create/Delete/etc. unintentionally
}
```

### 8. Use `slog` for Structured Logging (Go 1.21+)

```go
// ✅ GOOD
import "log/slog"

slog.Info("user created",
    slog.String("user_id", user.ID),
    slog.String("email",   user.Email),
)
slog.Error("payment failed",
    slog.String("order_id", orderID),
    slog.Any("error", err),
)

// ❌ BAD — unstructured, not queryable
log.Printf("user created: %s %s", user.ID, user.Email)
fmt.Println("payment failed:", err)
```

---

## Commands

```bash
go vet ./...           # static analysis
golangci-lint run      # comprehensive linting
go test ./...          # run all tests
go test -race ./...    # detect race conditions
go build ./...         # verify compilation
```

---

## Definition of Done

- [ ] All errors handled — no `underscore` for error returns
- [ ] Errors wrapped with `fmt.Errorf("context: %w", err)`
- [ ] `context.Context` passed as first parameter to all I/O functions
- [ ] `defer` used for all resource cleanup
- [ ] Interfaces defined at consumer side, kept to 1–3 methods
- [ ] Tests use table-driven pattern with `t.Run()`
- [ ] `go test -race ./...` passes with no race conditions
- [ ] `go vet ./...` passes clean

---

