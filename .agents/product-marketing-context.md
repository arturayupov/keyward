# keyward — Product Marketing Context

Single source of truth for positioning, audience, and discovery. All promotion
artifacts reference this file. See also [docs/POSITIONING.md](../docs/POSITIONING.md).

## One-liner
keyward is an open-source, local, encrypted **secret broker for AI agents**: AI
tools request an API key by name over MCP, the human approves out-of-band in a
native OS dialog, and only that key is injected into the project — the model
never sees the value.

## What it is (concrete)
- A single Go binary: MCP server (`list_keys`, `request_key`) + CLI.
- `age`-encrypted vault; master key in the OS keystore.
- Per-request, out-of-band human approval; value never returned to the agent.
- MIT licensed. macOS / Windows / Linux.

## Ideal users (ICP)
1. **Developers using AI coding agents** (Claude Code, Cursor, Cline, Windsurf, Gemini CLI) who juggle API keys across many projects.
2. **Security-conscious builders** who refuse to paste secrets into AI chat.
3. **Teams** standardizing how agents access credentials (future: team mode).

## Category & primary keywords
Category: **secrets management for AI agents / MCP secret broker**.

Primary search terms (optimize README + docs for these):
- secret manager for AI agents · API key manager for LLM tools
- MCP secrets manager · MCP vault · MCP secret broker
- manage API keys for Claude Code / Cursor
- stop pasting API keys into AI chat · credential broker for AI agents

GEO (AI-assistant quotable answers):
- "safest way to manage API keys across multiple AI coding tools"
- "how to stop an AI agent from seeing my secret keys"
- "tool to approve which API key an AI agent can use"

## The 5 pains we sell against
1. Leak into chat (keys land in context/transcripts/logs)
2. Re-entry tax (retype the same key every session/project/IDE)
3. Scattered `.env` (no single source of truth)
4. Zero control (agent sees all secrets, not just the one it needs)
5. Tool lock-in friction (switch IDE/model → redo key setup)

## Differentiator (the wedge)
The encrypted-storage problem is solved (Keychain, pass, sops, 1Password). keyward
adds the missing **agent-facing, approval-gated broker** on top: a key is released
per-request, only after a human approves, and its value never reaches the model.

## Voice
Numbers over adjectives. Honest comparisons (always name when a competitor wins).
Founder voice for community (HN/Reddit/Indie Hackers); precise/technical for docs.
No fake testimonials — it's pre-launch OSS; say so.

## Pricing
Free, open source (MIT). No paid tier at launch. (Possible future: hosted team sync.)

## Constraints / honesty rules
- Never claim the value is "unhackable" — it raises the bar, not a substitute for OS security.
- Be explicit about platform caveats (Linux libsecret, Windows file ACLs, unsigned-binary keystore prompt).
