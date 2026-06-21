# keyward — Design Spec

**Date:** 2026-06-21
**Status:** Approved for implementation planning
**Name:** `keyward` (decided — brandable, free on GitHub + npm; SEO via description/topics/README)

---

## 1. Problem

AI coding agents (Claude Code, Cursor, Gemini CLI, Windsurf, Cline, …) constantly need API keys and tokens — and today the human fumbles: hunting through scattered `.env` files, pasting secrets into chat (where they leak into context and transcripts), or re-entering the same key across every new session, project, IDE, and machine.

Existing secret managers (envchain, pass, sops, Infisical, 1Password CLI) solve *encrypted storage*. **None of them provide an agent-facing broker with per-request, out-of-band human approval.** That gap is what keyward fills.

## 2. One-line

A local, encrypted, **agent-facing secret broker**: any AI tool requests a *named* secret; the human approves a single request out-of-band via a native OS prompt; the broker **injects only that one key into the target** — the model never sees the value.

## 3. Goals / Non-goals

**Goals (v0):**
- The model **never receives a secret value** — broker injects into a target file/env and returns only a confirmation.
- Approval is **out-of-band** (native OS dialog), so the requesting agent cannot self-approve.
- **Universal**: works across any MCP-capable AI tool, plus a CLI for the rest; runs on macOS, Windows, Linux from a single binary.
- Secrets grouped by **project / direction** (namespaces), mirroring how the user already thinks about them.
- Import existing scattered `.env` files into the encrypted vault in one command.

**Non-goals (v0):**
- No cloud, no team/multi-user mode, no server component.
- No homegrown cryptography (use `age` + OS keystore).
- No biometric requirement in v0 (native confirm dialog suffices; biometric is a v1 hardening).
- No multi-device sync in v0 (explicitly deferred).

## 4. Architecture

```
   AI tool (Claude Code / Cursor / Gemini CLI / …)
        │  request_key("SHOPIFY_API_KEY", ns="shopify/livo", target="./.env")
        ▼
   ┌────────────────────────────────────────────┐
   │  keyward  (single Go binary)              │
   │                                             │
   │  ① Vault    age-encrypted file;             │
   │             master key in OS keystore       │
   │  ② Broker   policy check → approval → inject │
   │  ③ Approval native OS dialog (out-of-band)  │
   │  ④ Audit    append-only jsonl               │
   └────────────────────────────────────────────┘
        │  on approve: write KEY=val → target
        ▼  return ONLY "injected SHOPIFY_API_KEY → ./.env"
   AI tool continues — never saw the value
```

Three concerns are deliberately separated so each can vary independently (and so new OSes / tools plug in without touching the core):
- **Agent face** — how a tool asks (MCP server **or** CLI).
- **Human face** — how the user approves (native dialog per-OS, terminal fallback).
- **Vault** — how secrets are stored (age file + OS keystore).

## 5. Components & data model

### Storage layout (`~/.keyward/`)
| File | Purpose |
|---|---|
| `vault.age` | age-encrypted secrets blob; decrypted only in memory |
| `config.toml` | non-secret config (default targets, dialog backend, policy) |
| `audit.jsonl` | append-only decision log |

The vault's **age identity (master key)** is stored in the OS keystore via `go-keyring`:
- macOS → Keychain
- Windows → Credential Manager
- Linux → Secret Service (libsecret) / fallback to a passphrase-derived key

### Secret record
```
Secret {
  name           string      // e.g. "SHOPIFY_API_KEY"
  value          string      // never serialized to logs / stdout / MCP results
  namespace      string      // "shopify/livo", "ai/gemini", ...
  tags           []string
  allowed_targets []string   // optional glob allowlist, e.g. ["**/.env", "env"]
  created_at      time
  last_used_at    time
}
```

`namespace` is the project/direction grouping. A key may exist once and be referenced from multiple namespaces via tags, or be duplicated per project — import preserves whatever the source `.env` files imply.

## 6. Interfaces

### 6.1 MCP server (`keyward serve-mcp`)
Exposed tools — **values never appear in any tool result**:

- `list_keys(namespace?: string, query?: string)` → `[{ name, namespace, tags }]`
  Names and namespaces only. Lets the agent discover *what exists* without seeing secrets.
- `request_key(name, namespace, target, reason)` →
  - `target`: `"env"` (inject into the calling tool's process environment, where the transport supports it) **or** a filesystem path (e.g. `"./my-agent/.env"`).
  - Flow: resolve → policy check → **fire approval dialog** → on approve, inject → return `{ status: "injected", name, target }`. On deny/timeout: `{ status: "denied" }`.
  - **Return payload never contains the value.**

(Optional v0 helper, include if cheap: `suggest_keys(context)` — fuzzy-match a free-text need like "I need the Shopify token for livo" to candidate `{name, namespace}` pairs, still values-free.)

### 6.2 CLI
```
keyward init                       # create vault, generate+store master key
keyward import [path]              # scan .env files under path → vault, grouped by namespace
keyward add NAME --ns NS           # add a secret (value via prompt or stdin, never argv)
keyward ls [--ns NS]               # list names + namespaces (no values)
keyward inject NAME --ns NS --into ./.env   # human-driven inject (still triggers approval unless --yes in a TTY)
keyward get NAME --ns NS           # print value to stdout — gated, requires interactive confirm; for human use only
keyward audit [--tail N]           # show decision log
keyward serve-mcp                  # run the MCP server (stdio)
keyward watch                      # optional: foreground approval queue for headless/daemon use
```

## 7. Approval & security model (the keystone)

- The approval dialog is fired by the **keyward process**, never rendered into the agent's token stream. The agent therefore **cannot** click it or read its result.
- Dialog content: requesting tool (best-effort identification), key name, namespace, **target path**, reason string. Buttons: **Approve once · Approve for session · Deny**.
- "Approve for session" caches an approval for `{tool, namespace}` for a configurable TTL (default 30 min) held **in memory only**.
- v0 native dialog backends (built-in, no heavy deps):
  - macOS → `osascript -e 'display dialog …'`
  - Windows → PowerShell `System.Windows.MessageBox`
  - Linux → `zenity` / `kdialog`, else terminal prompt
  - Headless / no GUI → `keyward watch` terminal queue, or fail-closed.
- **Invariants (must hold, will be tested):**
  1. A secret value is never returned by any MCP tool.
  2. A secret value is never written to `audit.jsonl` or any log.
  3. `inject`/`request_key` never echo the value to stdout/stderr.
  4. Secret values are passed via prompt/stdin/file, never via process argv.
  5. Denied/timed-out requests inject nothing and are audited.

## 8. Crypto

- Vault encrypted with **`age`** (X25519, `filippo.io/age`). No custom crypto.
- Master age identity stored in the OS keystore (`github.com/zalando/go-keyring`). On Linux without a Secret Service, fall back to a passphrase-derived identity (scrypt).
- Vault is decrypted into memory on demand and zeroed after use where the language allows.

## 9. Part 1 integration — importing the user's existing keys

`keyward import ~/` (or a given root):
1. Walk for `.env` / `*.env` files, skipping `node_modules`, IDE/extension dirs, archives, and `.example`/`.template` files.
2. Parse `KEY=VALUE`; derive `namespace` from the project path (e.g. `~/livostyle/sync_config.env` → `livostyle`).
3. Load into the vault. On name collisions across projects, keep per-namespace copies (no silent merge).
4. Print a summary (names + namespaces only) and **offer** to delete or `chmod 600` the loose source files afterward — never automatic.

This satisfies the user's "collect all my keys into one protected place, grouped by project" request **without** ever creating a plaintext aggregate, and using the authoritative on-disk values rather than scraping session transcripts.

A pre-built inventory of the user's current keys (names + namespaces, no values) already exists at `~/.claude/secrets-inventory.md` and can seed the import's namespace mapping.

## 10. Open questions

- ~~**Final name.**~~ Resolved: **`keyward`** (free on GitHub + npm; `agentvault`/`keyper` were taken or collided in-category). Positioning, target queries, and pains captured in `docs/POSITIONING.md`.
- **`target: "env"` mechanics.** Injecting into the *calling tool's* process env is transport-dependent; for stdio MCP the realistic v0 target is a file path the tool then sources. Confirm during planning whether `env` target is feasible per client or should be deferred.
- **Tool identity.** How reliably can the MCP server identify the *requesting* tool for the dialog? v0: best-effort from MCP client info; never a security control, only a display hint.

## 11. Roadmap beyond v0

- v1: biometric approval (Touch ID / Windows Hello), tray/menubar app, per-key policy rules, rotation reminders.
- v2: multi-device sync (encrypted, user-owned backend), team mode.
- Distribution: Homebrew, Scoop, `go install`, GitHub releases. License: MIT or Apache-2.0.
