# Social posts (ready to paste)

Natural, founder-voice copy with SEO keywords woven in (MCP, API keys, AI agents,
Claude Code, Cursor, secrets management, secret broker, open source). Pick one per
platform. Replace the link if you set up a custom domain later.

Repo: https://github.com/arturayupov/keyward

---

## X / Twitter

### Option A — single post (tight)
Stop pasting API keys into your AI agent.

keyward lets Claude Code or Cursor request a key *by name* over MCP — you approve
it in a native OS popup, and only that key is injected into your project. The
model never sees the value.

Local. Encrypted. Open source.
https://github.com/arturayupov/keyward

### Option B — thread
1/ Every time I gave Cursor or Claude Code an API key, I pasted it straight into
the chat. Which means it's now in the context window, the logs, and the
transcript. Not great for a live Stripe key.

2/ Storing secrets is already solved — Keychain, 1Password, sops, pass. The part
nobody fixed is what happens when an *AI agent* needs one.

3/ So I built keyward. Your agent requests a key by name over MCP (`request_key`).
You get a native popup — approve once, approve for the session, or deny. On
approve, only that one key is injected into your project's .env. The model never
receives the value.

4/ Single Go binary. Local-first, age-encrypted vault, master key in your OS
keychain. MIT licensed, works with any MCP client — Claude Code, Cursor, Cline,
Windsurf.

brew install arturayupov/tap/keyward
https://github.com/arturayupov/keyward

---

## LinkedIn

I almost pasted a live Stripe key into an AI chat last week. Again.

If you use Claude Code, Cursor, or any AI coding agent, you know the moment: the
agent needs an API key, so you paste it into the conversation. Done — except that
secret now lives in the model's context, your logs, and the session transcript,
permanently.

Storing secrets is a solved problem. Keychain, 1Password, pass, sops — take your
pick. What nobody had solved is the handoff: how an AI agent uses a key without
ever seeing it.

So I built keyward, and open-sourced it today.

How it works:
• Your AI tool requests a key by name over MCP — the name, never the value.
• You get a native OS approval prompt: approve once, approve for the session, or deny.
• keyward injects only that key into your project's .env. The model never receives the value.

It's a single Go binary — local-first, encrypted at rest with age, master key in
your OS keystore. MIT licensed, already works with Claude Code, Cursor, Cline and
Windsurf, and it's listed in the official MCP registry.

If you've ever re-typed the same key across five projects, or winced right after
pasting one into a chat, take a look — and tell me where the approval flow feels
wrong. It's early and I'm listening.

https://github.com/arturayupov/keyward

#MCP #AIagents #DeveloperTools #SecretsManagement #OpenSource

---

## Notes
- **X:** keep hashtags to 0–1 (the culture dislikes hashtag spam); the thread
  works better than the single post for reach. Reply to early comments fast.
- **LinkedIn:** the first line is the hook (it's all that shows before "see more"),
  hashtags at the end are fine and help discovery.
- Both lead with the *pain* (pasting keys into chat), not the product name —
  that's what makes people stop scrolling.
