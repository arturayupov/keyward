# Contributing to keyward

Thanks for your interest! keyward is a small, focused Go project.

## Setup

```bash
git clone https://github.com/arturayupov/keyward
cd keyward
go test ./...    # all packages should pass (uses an in-memory keyring mock)
```

Requires Go 1.22+.

## Workflow

- **Tests first.** Every change ships with a test. Run `make test` (`go test ./... -race`).
- **Keep the security invariants.** The tests in `internal/vault`, `internal/audit`,
  `internal/broker`, `internal/mcp`, and `cmd/keyward` assert that secret values
  never leak (into logs, `ls`, MCP results, or `Result`). **Never weaken these.**
- **Format & vet.** `gofmt -w .` and `go vet ./...` must be clean.
- **Cross-platform.** Approval dialogs are per-OS behind build tags
  (`internal/approval/dialog_*.go`); verify `make cross` still builds all three.
  To verify a dialog **renders and parses correctly on your OS**, run the manual
  GUI test (excluded from CI):
  ```bash
  go test -tags manual -run TestLiveApprovalDialog ./internal/approval/ -v
  ```
  It pops the real native dialog; Approve must inject into the target, Deny must
  write nothing. Verified on macOS — Windows/Linux confirmation welcome.
- Small, focused PRs with a clear description. One concern per PR.

## Architecture

See [docs/specs/](docs/specs/) for the design and [docs/plans/](docs/plans/) for
the implementation plan. Each `internal/*` package has one responsibility:

| Package | Responsibility |
|---|---|
| `model` | Secret + value-free `SecretMeta` + `Store` |
| `vault` | age encryption + OS-keystore master key |
| `envfile` | parse + idempotent `0600` env writes |
| `importer` | walk `.env` files into namespaced secrets |
| `audit` | append-only, value-free decision log |
| `approval` | `Approver` interface, session cache, per-OS dialogs |
| `broker` | request → approve → inject → audit |
| `mcp` | `list_keys` / `request_key` tools |

## Code of conduct

Be respectful and constructive. Assume good faith.
