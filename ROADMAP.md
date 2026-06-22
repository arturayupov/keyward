# keyward Roadmap

Direction, not a contract — order and scope may shift with feedback. Have a
strong opinion on priorities? Open an issue or 👍 an existing one.

Status legend: ✅ done · 🚧 in progress · 📋 planned · 🔬 researching

---

## v0.1 — MVP ✅ (current)
The core broker, end to end.

- ✅ `age`-encrypted vault; master key in the OS keystore
- ✅ MCP server (`list_keys`, `request_key`) + CLI (`init`/`import`/`ls`/`inject`/`serve-mcp`)
- ✅ Out-of-band native approval dialogs (macOS / Windows / Linux), fail-closed
- ✅ Value-free audit log; "value never reaches the model" enforced by tests
- ✅ Import of scattered `.env` files (incl. `.env.local`/`.env.*`), namespaced by project
- ✅ Multi-line/special-char secret values (PEM keys, service-account JSON) inject intact
- ✅ Atomic vault writes; CI on macOS/Windows/Linux

## v0.2 — Distribution & frictionless install 📋 (next, post-launch)
Make it trivial to install and remove first-run friction.

- 📋 **Code signing + notarization** of release binaries — removes the first-run
  OS-keystore access prompt for end users. _Groundwork: `.goreleaser.yaml` is
  already signing-ready; needs Apple Developer ID + Windows cert._
- 📋 **Homebrew tap** (`brew install arturayupov/tap/keyward`) and **Scoop manifest**
- 📋 **Biometric approval** — Touch ID on macOS, Windows Hello — as a stronger,
  faster approval step than a button dialog
- 📋 **Windows ACL hardening** — set restrictive ACLs on injected files and the
  vault on Windows (Unix `0600` mode bits aren't honored there; today this is a
  documented limitation in [SECURITY.md](SECURITY.md))

## v1.0 — Control & ergonomics 📋
From "works" to "lives in your menubar."

- 📋 **Tray / menubar app** — see pending requests, approve/deny, recent activity
- 📋 **Per-key policy** — allowlist of targets/tools a key may be released to,
  optional auto-approve rules for trusted `{tool, namespace}` pairs
- 📋 **Key rotation reminders** — surface stale keys (uses the existing
  `LastUsedAt` field) and nudge rotation
- 📋 **`target: "env"` injection** — inject into a process's environment instead
  of a file, for tools that read secrets from `env` (deferred from the v0 spec)
- 🔬 **Secret-reference indirection** — let an agent reference a key it can use
  without even learning the name, for tighter least-privilege

## v2.0 — Multi-device & teams 🔬
- 🔬 **Encrypted multi-device sync** — user-owned backend (your storage, your
  keys); the sync layer never sees plaintext
- 🔬 **Team mode** — shared namespaces with per-member policy and approval routing
- 🔬 **Passphrase-derived fallback** for keystore-less environments (headless
  Linux/CI), as an alternative to a Secret Service provider

---

## How priorities are decided
1. **Don't weaken the core invariant** — the value never reaches the model. Any
   feature that would compromise this is off the table.
2. **Reduce friction before adding surface** — signing/install before fancy UI.
3. **Demand-driven** — issues with the most 👍 move up.

See [CONTRIBUTING.md](CONTRIBUTING.md) to help build any of these.
