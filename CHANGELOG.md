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

### Security
- Automated invariant tests asserting secret values never appear in logs,
  `ls` output, MCP results, or the broker `Result`.
