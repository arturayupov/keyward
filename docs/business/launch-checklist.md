# keyward — Launch & Submission Checklist

Paste-ready copy. **Everything here is GATED until the repo is public** and (for
PH/HN) ideally has a tagged release + a few GitHub stars for social proof. Order:
awesome-lists & registries first (in-category traffic), then HN/Reddit, then
Product Hunt.

---

## 0. Pre-flight (do before any submission)
- [ ] Repo public, README cover renders, CI badge green
- [ ] Tagged release (`v0.1.0`) with binaries (goreleaser)
- [ ] `go install ...@latest` verified from a clean machine
- [ ] Social-preview image set (Settings → Social preview → docs/images/cover.png)

---

## 1. awesome-mcp lists & MCP registries (highest in-category leverage)
Targets: `punkpeye/awesome-mcp-servers`, `modelcontextprotocol/servers` community
list, the official MCP registry, `wong2/awesome-mcp-servers`.

**One-line entry:**
> **[keyward](https://github.com/arturayupov/keyward)** — Local, encrypted secret broker. Agents request an API key by name; you approve in a native OS dialog; only that key is injected — the model never sees the value.

Category: *Security / Secrets*.

## 2. awesome-go / Go ecosystem
Targets: `avelino/awesome-go` (Security section), r/golang "what are you working on".

**Entry:**
> **[keyward](https://github.com/arturayupov/keyward)** — Encrypted, approval-gated secret broker for AI agents (MCP server + CLI).

## 3. Software directories
AlternativeTo (alternative to: "pasting API keys into AI chat", 1Password, envchain),
LibHunt (Go), SaaSHub, OpenAlternative (open-source directory).

**Short description (≤ 200 chars):**
> Open-source local secret broker for AI agents. Tools request API keys by name over MCP; you approve per-use; keys are injected into your project, never shown to the model.

---

## 4. Show HN (title + body)
**Title:** `Show HN: keyward – let AI agents use your API keys without ever seeing them`

**Body:**
> I kept pasting API keys into Claude Code/Cursor, which dumps them into context and transcripts — and re-typing the same keys in every new project. keyward is a small Go tool that fixes this: your keys live in one `age`-encrypted vault (master key in the OS keystore), and AI tools request a key **by name** over MCP. You approve the single request in a native OS dialog, and keyward injects only that key into the project's `.env`. The model never receives the value; `list_keys` returns names only, `request_key` returns a confirmation.
>
> It's MIT, single binary, macOS/Windows/Linux. Storage is a solved problem (Keychain, pass, sops, 1Password) — the missing piece was an agent-facing, approval-gated broker. Honest about limits: it raises the bar, it's not a substitute for OS security; Linux needs libsecret; Windows file ACLs are a v1 item.
>
> Repo + demo: <link>. Feedback welcome, especially on the approval UX and the threat model.

(Reply fast, founder voice, link the SECURITY.md when threat-model questions come.)

## 5. Reddit
**r/ClaudeAI / r/cursor — title:** `I built an open-source way to give Claude Code/Cursor your API keys without exposing them`

**Body (founder voice, no marketing speak):**
> Pasting keys into the agent always bugged me — they end up in context and logs. keyward keeps keys in an encrypted local vault and lets the agent request one by name over MCP; you approve in a native popup; only that key gets written into your `.env`, never shown to the model. MIT, single binary. Would love feedback on whether the approval flow feels right. <link>

**r/LocalLLaMA angle:** emphasize local-first, no cloud, you own the vault.

## 6. dev.to / blog post
**Title:** "Stop pasting API keys into your AI agent: a local approval-gated broker (MCP)"
Outline: the leak/re-entry problem → how MCP request_key + out-of-band approval
works → 60-second quickstart → threat model honesty → roadmap. Cross-post to
Hashnode; canonical on the repo if a site exists later.

## 7. Product Hunt (DRAFT only until traction)
**Tagline:** "Let AI agents use your API keys — without ever seeing them."
**First comment:** founder story (the pasting/re-entry pain), what it does in one
line, what it deliberately doesn't do, ask for the approval-UX feedback.
Do **not** launch PH until the repo has real usage/stars — PH rewards launch-day
social proof.

## 8. LinkedIn / X
Single post: the problem (one screenshot of a key pasted into chat → "this is in
your transcript forever"), the fix (gif of request → native approval → injected),
link. Tag MCP / Claude Code communities.
