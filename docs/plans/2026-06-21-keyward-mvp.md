# keyward MVP Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build keyward — a local, encrypted, agent-facing secret broker where AI tools request a named key, the human approves out-of-band, and the broker injects only that key into a target file; the model never sees the value.

**Architecture:** A single Go binary. Secrets live in an `age`-encrypted vault whose master identity is stored in the OS keystore. A `broker` orchestrates request → policy → approval → inject → audit. Two faces: a CLI (`cmd/keyward`) and an MCP server. Approval is an interface with per-OS native-dialog backends and a terminal fallback. Built in vertical slices: CLI vault+import+inject works before MCP and native dialogs are added.

**Tech Stack:** Go 1.22+, `filippo.io/age` (encryption), `github.com/zalando/go-keyring` (OS keystore), `github.com/mark3labs/mcp-go` (MCP server), `github.com/spf13/cobra` (CLI). Tests: stdlib `testing`.

**Conventions:** TDD (failing test first), DRY, YAGNI, commit after each green task, `gofmt -w .` before each commit. The security invariants from spec §7 are enforced by explicit tests in Tasks 6, 8, 10, 11, 13, 14 — those tests must never be weakened.

---

## File Structure

```
cmd/keyward/{main,init,import,ls,inject,serve}.go   # thin cobra commands
internal/config/config.go         # ~/.keyward paths
internal/model/secret.go          # Secret, SecretMeta, Store
internal/vault/{keystore,vault}.go# age identity in keystore; encrypt/decrypt
internal/envfile/envfile.go       # parse + idempotent KEY=VALUE write
internal/importer/importer.go     # walk .env files → namespaced secrets
internal/audit/audit.go           # append-only value-free jsonl
internal/approval/{approval,terminal,dialog_*}.go  # Approver + per-OS dialogs
internal/broker/broker.go         # request→approve→inject→audit
internal/mcp/server.go            # list_keys, request_key
```

Each `internal/*` package has one responsibility; `cmd/*` files just parse flags, call a package, print a redacted result. Full per-task code (tests + implementation) is captured in the task list below; this header is the map.

---

## Phases & Tasks (overview — each task is TDD: failing test → run red → implement → run green → commit)

### Phase 0 — Bootstrap
- **Task 1:** `go mod init github.com/arturayupov/keyward`; add deps (age, go-keyring, mcp-go, cobra); minimal cobra root in `cmd/keyward/main.go`; `go build ./...` passes.

### Phase 1 — Vault core
- **Task 2 (model):** `Secret{Name,Value,Namespace,Tags,AllowedTargets,CreatedAt,LastUsedAt}`, value-free `SecretMeta{Name,Namespace,Tags}`, `Store` with `Find/Upsert/Meta`. Tests: upsert-replaces-same-key; `Meta` carries no value.
- **Task 3 (config):** `PathsFor(home)` → `~/.keyward/{vault.age,audit.jsonl,config.toml}`; `EnsureDir` 0700. Test: path resolution.
- **Task 4 (keystore):** `EnsureIdentity()` gets/creates an age X25519 identity in the OS keystore (`go-keyring`); tests use `keyring.MockInit()`. Test: identity stable across calls.
- **Task 5 (vault):** `Save(path,store,recipient)` (age-encrypt JSON, 0600) / `Load(path,id)`. Tests: round-trip; **encrypted-at-rest** (plaintext value absent from file bytes).

### Phase 2 — env files & import
- **Task 6 (envfile) [SECURITY]:** `Parse`/`ParseFile` (handles `export`, comments, quotes); `Set(path,key,value)` idempotent, 0600, replaces-not-duplicates, never logs/returns value. Tests: parse; perms; replace-not-duplicate.
- **Task 7 (importer):** `Import(root)` walks `.env`/`*.env`, skips `node_modules`/IDE dirs/`.example`/`.template`, namespace = first path segment under root. Test: walks + namespaces + skips.

### Phase 3 — audit, approval, broker
- **Task 8 (audit) [SECURITY]:** `Entry` struct is value-free by construction; `Record` appends jsonl 0600. Test: name present, no `value` field.
- **Task 9 (approval):** `Approver` interface; `Decision{Deny,ApproveOnce,ApproveSession}`; `Request` (value-free); `SessionCache` (caches only ApproveSession per {tool,ns,name}); `TerminalApprover` fallback. Tests: session cache avoids re-prompt; ApproveOnce not cached.
- **Task 10 (broker) [SECURITY]:** `Broker.Request(req)` → find secret → approve → on approve `envfile.Set` else nothing → audit; `Result` is value-free. Tests: approved injects; denied injects nothing; result has no value field.

### Phase 4 — CLI (working tool end-to-end)
- **Task 11 (cli) [SECURITY]:** `init`, `import [root]`, `ls [--ns]`, `inject NAME --ns --into`. In-process cobra tests with temp `HOME` + `keyring.MockInit()` in `TestMain`. Test: `ls` prints names/namespaces, never a value. After this task keyward is a usable CLI.

### Phase 5 — Native dialogs
- **Task 12 (dialogs):** `nativeApprover()` behind build tags — `dialog_darwin.go` (osascript), `dialog_windows.go` (PowerShell MessageBox), `dialog_linux.go` (zenity, else false), `dialog_other.go` (false). Exported `NativeApprover()`; `inject`/`serve` pick native else terminal. Verify `GOOS=darwin|windows|linux go build ./...` all pass. Fail-closed on dialog error/cancel → Deny.

### Phase 6 — MCP server
- **Task 13 (mcp) [SECURITY]:** `Handlers{Store,Broker}` with `listKeys`/`requestKey`; `NewServer` registers `list_keys` (returns `Store.Meta` JSON) and `request_key` (returns `broker.Result` JSON — `injected … → target`, never the value). `serve-mcp` command runs `server.ServeStdio`. Tests: list returns no values; request returns confirmation not value. **Verify mcp-go API against current docs (context7 / pkg.go.dev) before coding** — registration shape may have drifted.

### Phase 7 — Hardening & release readiness
- **Task 14 (invariants + CI) [SECURITY]:** cross-package test: full `init→import→inject` with sentinel value `SENTINEL-d4f9`; assert it appears ONLY in `vault.age` (encrypted) and the injected target — never in stdout, `audit.jsonl`, or config. `.github/workflows/ci.yml`: matrix {ubuntu,macos,windows} × `go build/vet/test -race`.
- **Task 15 (release):** `.goreleaser.yaml` (darwin/linux/windows, amd64+arm64) + `Makefile` (`build/test/install`). No release published while private.

### Phase 8 — Polished OSS repo (needs working tool for screenshots)
- **Task 16 (README):** expand to full structure — hero+cover+badges, the 5 pains (from `docs/POSITIONING.md`), how-it-works diagram + value-never-seen guarantee, install (`go install`, brew, scoop, binaries), quickstart (`init`→`import ~/`→MCP config→agent request→approve), CLI table, security model, honest comparison vs envchain/pass/sops/1Password, roadmap/contributing/license. Include the exact Claude Code MCP config:
  ```json
  { "mcpServers": { "keyward": { "command": "keyward", "args": ["serve-mcp"] } } }
  ```
- **Task 17 (docs):** `INSTALL.md` (per-OS, incl. libsecret on Linux), `USAGE.md` (real session transcript), `SECURITY.md` (threat model + disclosure), `CONTRIBUTING.md`, `LICENSE` (decide MIT vs Apache-2.0 here), `CHANGELOG.md`.
- **Task 18 (visuals):** `docs/images/cover.png` (1280×640 social preview — name + "secret broker for AI agents" + flow glyph, via canvas-design skill), screenshots (native dialog mid-request, `ls` output, agent receiving "injected"), `quickstart.gif` (<10s). Wire into README; set GitHub social preview.
- **Task 19 (pre-public gate):** run `/code-review`; manual e2e per OS; verify no secret in git history (sentinel grep) + `.gitignore` covers all artifacts; **only after user sign-off** flip to public (`gh repo edit --visibility public`), publish release, submit to awesome-mcp lists.

---

## Detailed code reference

Full TDD code for Tasks 2–13 (tests + implementations) is provided inline during execution. The non-obvious, get-it-right-once pieces are pinned here so they are not re-derived:

**age vault (Task 5):** `age.Encrypt(w, recipient)` → write JSON → `Close()` → `os.WriteFile(path, …, 0600)`. Load: `age.Decrypt(f, id)` → `io.ReadAll` → `json.Unmarshal`. Encrypted-at-rest test greps raw file bytes for the plaintext and fails if found.

**keystore (Task 4):** service `"keyward"`, account `"vault-age-identity"`; `keyring.Get` → `age.ParseX25519Identity`; on `ErrNotFound` → `age.GenerateX25519Identity` + `keyring.Set(id.String())`.

**envfile.Set (Task 6):** read existing lines, replace the line whose trimmed `export?KEY=` matches (single occurrence), else append; quote values containing whitespace/quotes/`#`; write `0600`. Never include the value in any returned error.

**approval (Task 9):** `SessionCache.Approve` returns cached `ApproveSession` without calling inner; only `ApproveSession` is cached (key = `tool\x00ns\x00name`). All dialog backends **fail closed** (error/cancel → `Deny`).

**broker (Task 10):** `Result{Status,Name,Target}` — no value field, ever. Deny path writes nothing to the target and audits `"denied"`.

**mcp (Task 13):** handlers return JSON of `Store.Meta` / `broker.Result`; never marshal a `Secret`. `serve-mcp` loads the vault once, builds one `Broker` with `SessionCache(native||terminal)`.

---

## Self-Review (by plan author)

- **Spec coverage:** §4 arch → T2–13; §5 storage/model → T2–5,8; §6 MCP+CLI → T11,13; §7 approval + invariants → T6,8,9,10,12,14; §8 crypto → T4,5; §9 import → T7; §10 open questions (`target:env`, tool identity) → deferred, noted in T13 (v0 target is a file path); §11 roadmap/distribution → T15,18,19; README/positioning → Phase 8.
- **Placeholder scan:** no "implement later" in code steps; the two deferred items (biometric approval, process-`env` injection) are explicit roadmap, not gaps; mcp-go API drift is flagged with a concrete verify step.
- **Type consistency:** `Decision/Request/Approver/Result/Store.Meta/SecretMeta/envfile.Set/audit.Entry/broker.Broker` used identically across tasks.

## Open items to confirm during execution
1. **mcp-go API** — verify registration calls against current SDK (Task 13).
2. **License** — MIT vs Apache-2.0 (Task 17).
3. **`target:"env"`** process-env injection — deferred; v0 uses a file-path target.
4. **Disk space** — building Go needs the module cache + build cache (hundreds of MB). The machine is currently at ~100% (≈0.9 GB free); free space before Task 1 or builds will fail with ENOSPC.
