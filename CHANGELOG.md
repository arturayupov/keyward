# Changelog

All notable changes to keyward are documented here. The format is based on
[Keep a Changelog](https://keepachangelog.com/) and this project adheres to
[Semantic Versioning](https://semver.org/).

## [Unreleased]

### Added
- Encrypted vault (`age`) with the master key stored in the OS keystore
  (macOS Keychain / Windows Credential Manager / Linux libsecret).
- `keyward init` / `import` / `ls` / `inject` / `serve-mcp` CLI.
- MCP server tools `list_keys` and `request_key` (values never returned to the agent).
- Out-of-band native approval dialogs for macOS, Windows, and Linux, with a
  terminal fallback; all backends fail closed.
- Append-only, value-free audit log.
- Import of existing `.env` files, grouped into namespaces by project.
- Friendly guidance when the OS keystore is unreachable.
- CI matrix (ubuntu/macos/windows), goreleaser config, Makefile.

### Fixed
- Multi-line and special-character secret values (PEM keys, service-account
  JSON, values with quotes/backslashes/`#`) now round-trip through `.env`
  injection intact instead of corrupting the file.
- Vault writes are atomic (temp file + rename) — an interrupted write can no
  longer truncate or corrupt an existing vault.
- Import now also picks up `.env.local` / `.env.production` / `.env.*` files.

### Security
- Automated invariant tests asserting secret values never appear in logs,
  `ls` output, MCP results, or the broker `Result`.
- Concurrent `request_key` injections into the same target file are serialized.

### Docs
- [ROADMAP.md](ROADMAP.md) — versioned plan (v0.2 distribution, v1.0 control,
  v2.0 multi-device/teams).
