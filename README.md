# keyward — a secret broker for AI agents

**Store your API keys once. Let any AI tool request them with your approval — without ever seeing the value.**

keyward is an open-source, local, encrypted **secret broker for AI coding agents**. Instead of pasting API keys into chat (where they leak into context and transcripts) or re-entering the same key in every new session, project, and IDE, you keep your keys in one encrypted vault. When Claude Code, Cursor, Gemini CLI, or any [MCP](https://modelcontextprotocol.io)-capable tool needs a key, it **requests it by name**, you **approve that single request** in a native OS prompt, and keyward **injects only that key** into your project. The model never receives the value.

> Status: design complete, implementation starting. See [`docs/specs`](docs/specs/2026-06-21-keyward-design.md).

## Why

- **Stop leaking keys into AI chat.** Secrets are injected into your `.env`/process, never returned to the model.
- **One vault, every tool.** Works across any MCP client plus a CLI — same keys in Claude Code, Cursor, Gemini CLI, and more, on macOS, Windows, and Linux.
- **You stay in control.** Each request triggers an out-of-band approval the agent cannot click. Approve once, or for the session.
- **No more scattered `.env` files.** Import your existing keys into one encrypted, project-grouped store.

## How it works

```
AI tool ──request_key("STRIPE_KEY", project)──▶ keyward ──native approval──▶ you
                                                    │ approved
        ◀── "injected STRIPE_KEY → ./.env" ─────────┘   (value never shown to the model)
```

- **Vault:** [age](https://age-encryption.org)-encrypted file; master key in the OS keystore (Keychain / Windows Credential Manager / libsecret).
- **Interfaces:** an MCP server (`list_keys`, `request_key`) and a `keyward` CLI.
- **Approval:** native OS dialog, fired by keyward itself — out-of-band, so the agent can't self-approve.

## Install

_Coming soon — `go install`, Homebrew, Scoop. Single static binary, no runtime._

## License

TBD (MIT or Apache-2.0).
