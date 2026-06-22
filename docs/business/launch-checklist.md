# keyward — Launch Guide (step-by-step)

Plain-language guide: what each place is, the exact URL, where to click, and
what to paste. Nothing here can be auto-posted (every site needs *your* login,
and HN/Reddit/PH forbid bot posting) — but each item below is a 2-minute
copy-paste.

**Do these in order.** Top ones are highest-leverage and lowest-effort.

Reusable links:
- Repo: https://github.com/arturayupov/keyward
- Cover image: `docs/images/cover.png`
- One-liner: *Local, encrypted secret broker for AI agents — tools request an API key by name over MCP, you approve in a native OS dialog, and only that key is injected into your project; the model never sees the value.*

---

## ✅ Already done (by us)
- **Official MCP Registry** — published as `io.github.arturayupov/keyward`.
- **awesome-mcp-servers** — PR opened (#8509), waiting on maintainer merge.
- **Homebrew / Scoop / Go install / GitHub Releases** — all live.

## 🤖 Happens automatically now (no action)
Several MCP directories crawl the **official registry** and public GitHub repos
with the `mcp` topic, so they tend to list new servers within days:
- **Glama** (glama.ai/mcp/servers) — auto-indexes public MCP repos.
- **mcp.so**, **PulseMCP** (pulsemcp.com) — pull from the registry/GitHub.
Check back in ~1 week; if missing, each has a "Submit" button (see §5).

---

## 1. GitHub social preview (2 min) — DO FIRST
**What:** the preview image shown when your repo link is pasted into Twitter/X,
Slack, Discord, iMessage. Without it you get a generic gray box.
**How:** GitHub → your repo → **Settings** tab → scroll to **Social preview** →
**Edit → Upload an image** → pick `docs/images/cover.png`.

## 2. GitHub Release notes polish (2 min)
**What:** the v0.1.1 release page is the first thing visitors check.
**Where:** https://github.com/arturayupov/keyward/releases/tag/v0.1.1 → **Edit**.
**Paste** a short "what is this" intro above the auto-generated changelog (use
the one-liner + the install commands from the README).

## 3. Hacker News — "Show HN" (5 min) — biggest single spike
**What is HN:** news.ycombinator.com — the dev community where "Show HN" posts
present things you built. A front-page Show HN can bring thousands of devs.
**Rules:** post it yourself, plain link, no marketing speak; reply to every
comment fast; post once (no reposting). Best time: weekday ~8–10am US Eastern.
**Where:** https://news.ycombinator.com/submit (must be logged in).
- **Title:** `Show HN: keyward – let AI agents use your API keys without ever seeing them`
- **URL:** https://github.com/arturayupov/keyward
- **First comment (paste right after posting):**
> I kept pasting API keys into Claude Code/Cursor, which dumps them into context and transcripts — and re-typing the same keys in every project. keyward keeps keys in one age-encrypted vault (master key in the OS keystore); AI tools request a key **by name** over MCP, you approve the single request in a native OS dialog, and only that key is injected into the project's `.env`. The model never receives the value. MIT, single Go binary, macOS/Windows/Linux. Storage is solved (Keychain, pass, sops, 1Password) — the missing piece was an agent-facing, approval-gated broker. Honest limits: it raises the bar, not a substitute for OS security; Linux needs libsecret; Windows file ACLs are a v1 item. Feedback welcome, especially on the approval UX.

## 4. Reddit (5 min each) — targeted communities
**What is Reddit:** topic forums ("subreddits"). Post as yourself; most subs
**dislike pure self-promotion** — frame it as "I built this to solve X, feedback?"
and engage in comments. Read each sub's rules first (right sidebar).
**Where / which subs:**
- **r/mcp** — https://www.reddit.com/r/mcp/submit — the MCP community (most on-topic).
- **r/ClaudeAI** — https://www.reddit.com/r/ClaudeAI/submit
- **r/LocalLLaMA** — https://www.reddit.com/r/LocalLLaMA/submit — emphasize local-first / you own the vault.
- **Title:** `I built an open-source way to give Claude Code/Cursor your API keys without exposing them`
- **Body:**
> Pasting keys into the agent always bugged me — they end up in context and logs. keyward keeps keys in an encrypted local vault and lets the agent request one by name over MCP; you approve in a native popup; only that key gets written into your `.env`, never shown to the model. MIT, single binary. Would love feedback on whether the approval flow feels right. https://github.com/arturayupov/keyward

## 5. MCP directories — manual submit (3 min each)
Only needed if they haven't auto-listed you after ~1 week.
- **Smithery** — https://smithery.ai → sign in with GitHub → **Add Server** / "Deploy" → point at the repo. (Largest MCP marketplace.)
- **Glama** — https://glama.ai/mcp/servers → usually auto-added; if not, there's a claim/submit flow after GitHub login.
- **mcp.so** — https://mcp.so → "Submit" (GitHub URL).
- **PulseMCP** — https://www.pulsemcp.com → "Submit a server".
- **Cursor Directory** — https://cursor.directory/mcp → submit (GitHub-based).
- **Cline Marketplace** — open a PR/issue at https://github.com/cline/mcp-marketplace per their template.

## 6. dev.to — a launch blog post (20 min) — good for SEO
**What is dev.to:** a developer blogging platform; posts rank in Google and get
shared. Free account.
**Where:** https://dev.to/new (after signing up).
- **Title:** "Stop pasting API keys into your AI agent: a local approval-gated broker (MCP)"
- **Outline:** the leak/re-entry problem → how MCP `request_key` + out-of-band
  approval works → 60-second quickstart → threat-model honesty → roadmap.
- Add tags: `mcp`, `ai`, `security`, `go`. Cross-post to Hashnode if you want.

## 7. Software directories (3 min each) — long-tail discovery
**What:** sites where people search for/compare tools.
- **AlternativeTo** — https://alternativeto.net → "Add application" → list keyward as an alternative to "pasting API keys into AI chat", envchain, 1Password. Description (≤200 chars): *Open-source local secret broker for AI agents. Tools request API keys by name over MCP; you approve per-use; keys are injected, never shown to the model.*
- **LibHunt (Go)** — https://www.libhunt.com → submit the GitHub repo.
- **OpenAlternative** — https://openalternative.co → "Submit" (open-source tools).
- **SaaSHub** — https://www.saashub.com → "Submit a product".

## 8. Product Hunt (DRAFT now, launch later) — wait for traction
**What is PH:** producthunt.com — a daily launch board; launch day rewards
upvotes/comments, so **do it only once you have GitHub stars / a few users** as
social proof. Prepare the draft now, launch in a few weeks.
**Where:** https://www.producthunt.com/products/new (sign in).
- **Tagline (≤60 chars):** "Let AI agents use your API keys — without ever seeing them."
- **First comment:** the founder story (the pasting/re-entry pain), one line on
  what it does, what it deliberately doesn't do, ask for approval-UX feedback.

## 9. X / LinkedIn (2 min) — your own audience
**Post idea:** screenshot of a key pasted into chat ("this is in your transcript
forever") → short gif of request → native approval → injected → link. Tag the
MCP / Claude Code communities.

---

## Etiquette that keeps you out of trouble
- **One post per platform.** Reposting reads as spam.
- **Founder voice** on HN/Reddit/IndieHackers; vendor voice on directories.
- **Reply fast** the first few hours — that's what drives ranking.
- **Be honest about limits** (it's pre-1.0, unsigned binaries prompt the keystore,
  Linux needs libsecret). Devs trust honesty and punish hype.
- **No fake stars/upvotes.** PH and HN detect and penalize it.
