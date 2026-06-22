# keyward — Promotion Status

State tracker for the promotion playbook, adapted for a **pre-public OSS repo
with no website/domain**. Items needing a live site or a public repo are marked
🔒 GATED. Update checkboxes as work lands.

## Phase 0 — Context + state ✅
- [x] `.agents/product-marketing-context.md`
- [x] `.agents/promotion-status.md`

## Phase 1 — Discovery foundation (repo-level)
- [x] GitHub description set
- [x] GitHub topics set (17)
- [x] `llms.txt` at repo root (AI-citation surface)
- [x] `README.md` optimized for primary keywords (headings/body)
- [x] `SECURITY.md` (doubles as security disclosure surface)
- [ ] 🔒 GATED: JSON-LD schema, sitemap, robots, IndexNow, security.txt — require a website/domain
- [ ] GitHub social-preview upload (manual; Settings → Social preview → docs/images/cover.png)

## Phase 2 — AI-discovery layer
- [x] MCP server (keyward *is* one — `serve-mcp`)
- [ ] 🔒 GATED (post-public): Wikidata entity (needs notability/public presence)
- [ ] 🔒 GATED (post-public): Custom GPT / Claude Project / Perplexity Collection loaded with llms.txt
- [ ] Submit to MCP registries / `awesome-mcp` (copy ready in launch-checklist) — needs public repo

## Phase 3 — Comparison content
- [x] README comparison table (envchain / pass-sops / 1Password)
- [x] `docs/comparisons/*.md` long-form honest comparisons (with "when X wins")

## Phase 4 — Directory / list submissions (copy ready, submission GATED on public)
- [x] Paste-ready copy written (`docs/business/launch-checklist.md`)
- [x] awesome-mcp PR opened (punkpeye/awesome-mcp-servers#8509); MCP registry server.json ready (publish = user `mcp-publisher login github`)
- [ ] 🔒 AlternativeTo, LibHunt, SaaSHub
- [ ] 🔒 Product Hunt (draft only until traction)

## Phase 5 — Community + content (drafts ready, posting GATED on launch)
- [x] Launch posts drafted (`docs/business/launch-checklist.md`)
- [ ] 🔒 Show HN, r/ClaudeAI, r/LocalLLaMA, dev.to, LinkedIn

## Phase 6 — Brand mark
- [x] Cover image (`docs/images/cover.png`)
- [ ] N/A: web favicon (no website yet)

## Pre-public gate (must pass before flipping to public)
- [ ] `/code-review` of the repo
- [ ] Live MCP-in-Claude-Code test (per user: after publication)
- [ ] Final secret/history scan (done once; re-run before flip)
- [ ] User sign-off → `gh repo edit --visibility public` → tag release
