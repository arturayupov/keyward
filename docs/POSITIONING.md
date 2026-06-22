# keyward — Positioning & Discovery

**Date:** 2026-06-21

## Name
**keyward** — "key + ward (guardian)". Free on GitHub and npm. Brandable; SEO/GEO is carried by the description, topics, and README — not the name itself (no new project ranks on an unknown brand word).

## GitHub "About" description (≤ 350 chars)
> Open-source secret broker for AI agents. Store your API keys once; let Claude Code, Cursor, and any MCP tool request them with per-use human approval — keys are injected into your project, never exposed to the model.

## GitHub topics
`mcp` · `model-context-protocol` · `secrets-management` · `api-keys` · `ai-agents` · `claude-code` · `cursor` · `llm` · `vault` · `credentials` · `developer-tools` · `cli` · `golang` · `security` · `dotenv`

## Target search queries (optimize README headings/body for these)
**Search engines (Google/Bing):**
- secret manager for AI agents
- API key manager for LLM / AI tools
- MCP secrets manager · MCP vault · MCP secret broker
- manage API keys for Claude Code / Cursor
- single place to store API keys for AI
- stop pasting API keys into AI chat
- credential broker for AI agents

**AI assistants (GEO — README should be the quotable answer):**
- "What's the safest way to manage API keys across multiple AI coding tools?"
- "How do I stop my AI agent from seeing my secret keys?"
- "Is there a tool to approve which API key an AI agent can use?"

## Pains we sell against
1. **Leak into chat** — pasting a key into a prompt lands it in context, transcripts, logs.
2. **Re-entry tax** — same key typed again every new session / project / IDE / machine.
3. **Scattered `.env`** — no single source of truth; you forget which key lives where.
4. **Zero control** — an agent can see/grab every secret at once instead of the one it needs.
5. **Tool lock-in friction** — switching IDE or model means redoing all key setup.

## Distribution checklist (post-MVP)
- Submit to `awesome-mcp` lists and MCP registries (largest single source of in-category traffic).
- `go install`, Homebrew tap, Scoop manifest, GitHub Releases binaries.
- Launch posts: Product Hunt, Hacker News (Show HN), dev.to / r/LocalLLaMA / r/ClaudeAI.
- License: MIT or Apache-2.0 (decide at first public push).
