# AGENTS.md — go-email-normalizer

Guidance for agentic coding tools (Copilot, Cursor, OpenCode, etc.) operating in this repository.

---

## Repository Overview

`go-email-normalizer` is a Go library (not a binary) that normalizes email addresses across major providers (Apple, Fastmail, Google, Microsoft, Protonmail, Rackspace, Rambler, Yahoo, Yandex, Zoho).

- **Module:** `github.com/bobadilla-tech/go-email-normalizer/v5`
- **Minimum Go version:** 1.14
- **Package name:** `emailnormalizer` (all files in root, no subdirectories)
- **Only external dependency:** `github.com/stretchr/testify v1.6.1` (test assertions)

---

## Build, Test, and Lint Commands

### Run all tests
```sh
go test ./...
```

### Run all tests with verbose output (mirrors CI)
```sh
go test -v ./...
```

### Run all tests with coverage (exactly as CI does)
```sh
go test -v -coverprofile=coverage.txt ./...
```

### Run a single test function
```sh
go test -v -run TestGoogleUsername ./...
go test -v -run TestNormalizer_InvalidEmail ./...
```

### Run a single test file's tests (by matching a prefix)
```sh
go test -v -run TestGoogle ./...
go test -v -run TestYahoo ./...
```

### Build (verify compilation — this is a library, no binary)
```sh
go build ./...
```

### Vet (static analysis)
```sh
go vet ./...
```

### Format
```sh
gofmt -w .
# or
goimports -w .
```

There is no Makefile and no `.golangci.yml`. No linter is configured in CI — only `go test` and `go vet` are expected to pass.

---

## Project Structure

The repository is a **flat single-package layout**. All `.go` files live in the root directory.

```
go-email-normalizer/
├── rule.go              # NormalizingRule interface
├── normalizer.go        # Normalizer struct, NewNormalizer, Normalize, AddRule
├── domains.go           # Per-provider domain slice vars
├── <provider>_rule.go   # One file per provider (implements NormalizingRule)
└── <provider>_rule_test.go
```

When adding a new provider, follow this pattern:
1. Add domain(s) to `domains.go` as a `var <provider>Domains = []string{...}`.
2. Create `<provider>_rule.go` with a struct implementing `NormalizingRule`.
3. Register all domains in `NewNormalizer()` inside `normalizer.go`.
4. Create `<provider>_rule_test.go` with `Test<Provider>Username` and `Test<Provider>Domain`.

---

## Code Style Guidelines

### Formatting
- All code must be formatted with `gofmt`. Never submit unformatted code.
- Use tabs for indentation (Go standard).
- No trailing whitespace.

### Imports
- Use the standard Go import grouping:
  1. Standard library
  2. Third-party packages
- Separate groups with a blank line.
- Only `strings` from stdlib and `github.com/stretchr/testify/assert` appear in this codebase; keep imports minimal.

### Naming Conventions
- **Package:** `emailnormalizer` (single lowercase word, no underscores).
- **Rule structs:** `<Provider>Rule` (e.g., `GoogleRule`, `YahooRule`).
- **Domain variables:** `<provider>Domains` (e.g., `googleDomains`, `yahooDomains`), unexported.
- **Test functions for provider rules:** `Test<Provider>Username` / `Test<Provider>Domain`.
- **Test functions for Normalizer scenarios:** `TestNormalizer_<Scenario>` (e.g., `TestNormalizer_InvalidEmail`).
- Use `camelCase` for unexported identifiers and `PascalCase` for exported ones.
- Avoid abbreviations unless they are universally understood (e.g., `n` as a receiver for `*Normalizer` is idiomatic).

### Types and Interfaces
- `NormalizingRule` is the core interface — every provider must implement it exactly:
  ```go
  type NormalizingRule interface {
      ProcessUsername(string) string
      ProcessDomain(string) string
  }
  ```
- Do not add methods to the interface without a compelling reason; it is intentionally minimal.
- Provider rule structs have no fields; they are stateless value types.

### Error Handling
- This library does not return errors. Invalid or unrecognizable emails are returned unchanged (the `Normalize` function returns the input as-is on parse failure).
- Do not introduce `error` return values to `Normalize`, `ProcessUsername`, or `ProcessDomain`.
- Edge cases (missing `@`, multiple `@`, empty string) are handled silently in `Normalize`.

### String Processing
- All username processing must begin with `strings.ToLower`.
- Domain is always lowercased in `Normalize` before the rule is applied; `ProcessDomain` receives an already-lowercase domain.
- Use `strings` stdlib functions. Do not use `regexp` unless a transformation cannot reasonably be expressed with `strings` functions.

---

## Testing Guidelines

### Framework
- Use `github.com/stretchr/testify/assert` for all assertions.
- All tests are in package `emailnormalizer` (not `emailnormalizer_test`) — white-box testing.

### Test Style
- **Preferred style:** flat test functions with multiple `assert.Equal(t, expected, actual)` calls.
- **Table-driven tests** (`t.Run`) are acceptable for complex input matrices (see `protonmail_rule_test.go` for an example).
- Do not mix both styles in the same test function.
- Keep tests in the same package as the code under test (no `_test` package suffix).

### Table-driven test template (when needed)
```go
func TestProviderRule_ProcessUsername(t *testing.T) {
    tests := []struct {
        name string
        args string
        want string
    }{
        {"plain", "User", "user"},
        {"plus tag", "user+tag", "user"},
    }
    rule := ProviderRule{}
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            assert.Equal(t, tt.want, rule.ProcessUsername(tt.args))
        })
    }
}
```

### Coverage
- CI uploads coverage to Codecov. Aim to cover both `ProcessUsername` and `ProcessDomain` for every rule, and the key edge cases in `Normalize`.

---

## CI/CD

The GitHub Actions workflow (`.github/workflows/go.yml`) triggers on every push and:
1. Checks out the repo.
2. Sets up Go (latest via `actions/setup-go@v5`).
3. Runs `go mod download`.
4. Runs `go test -v -coverprofile=coverage.txt ./...`.
5. Uploads coverage to Codecov.

**A PR is considered passing when `go test ./...` exits 0 and `go build ./...` succeeds.**

---

## Common Pitfalls

- Do not add `vendor/` — the project uses the module cache.
- Do not introduce build tags; all files compile unconditionally.
- `ProcessDomain` must return a non-empty string; returning an empty string will produce a malformed `@`-prefixed address.
- When adding a new provider, register **all** known domains in `NewNormalizer()`, not just the primary one. Alias domains (e.g., `googlemail.com`) must also be registered and optionally canonicalized in `ProcessDomain`.
