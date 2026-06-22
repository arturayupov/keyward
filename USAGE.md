# Using keyward

## 1. Initialize

```bash
$ keyward init
Initialized vault at /Users/you/.keyward/vault.age
```

This generates an `age` identity, stores it in your OS keystore, and writes an
empty encrypted vault.

## 2. Add keys

### Import from existing `.env` files
```bash
$ keyward import ~/projects
Imported 7 secrets (names only):
  shop-app/SHOPIFY_API_KEY
  shop-app/SHOPIFY_API_SECRET
  bot/TELEGRAM_BOT_TOKEN
  ...
```
Namespaces are derived from the project folder. `.env.example`/`.template` files,
`node_modules`, and IDE/cache directories are skipped.

### List what you have (names only — never values)
```bash
$ keyward ls
shop-app    SHOPIFY_API_KEY
shop-app    SHOPIFY_API_SECRET
bot         TELEGRAM_BOT_TOKEN
```

## 3. Use a key from an AI tool (MCP)

Configure the MCP server once per tool.

**Claude Code** — `~/.claude.json` or a project `.mcp.json`:
```json
{ "mcpServers": { "keyward": { "command": "keyward", "args": ["serve-mcp"] } } }
```
Cursor / Windsurf / Cline use the same `command`/`args`.

Now ask the agent naturally:

> "Add my SHOPIFY_API_KEY (shop-app) to this project's .env so the dev server can start."

What happens:
1. The agent calls `list_keys` to discover the name — it sees `SHOPIFY_API_KEY` / `shop-app`, **not the value**.
2. The agent calls `request_key(name, namespace, target=".env", reason=...)`.
3. **keyward shows a native dialog**: *"mcp requests "SHOPIFY_API_KEY" (project "shop-app") → ./.env"* with **Approve once / Approve for session / Deny**.
4. On approval, keyward writes `SHOPIFY_API_KEY=…` into `./.env` (`0600`) and returns to the agent only: `{"status":"injected","name":"SHOPIFY_API_KEY","target":"./.env"}`.
5. The agent continues. It never saw the secret.

"Approve for session" remembers that `{tool, namespace, key}` for the rest of the
process so repeated requests don't re-prompt.

## 4. Use a key from the CLI

```bash
$ keyward inject SHOPIFY_API_KEY --ns shop-app --into ./.env
# (approval prompt appears)
injected SHOPIFY_API_KEY → ./.env
```

## Where things live

| Path | What |
|---|---|
| `~/.keyward/vault.age` | age-encrypted secrets (`0600`) |
| OS keystore | the vault's master key |
| `~/.keyward/audit.jsonl` | append-only decision log (no values) |

## Audit

Every request — approved or denied — is recorded value-free:
```json
{"time":"2026-06-21T16:57:02Z","tool":"mcp","name":"SHOPIFY_API_KEY","namespace":"shop-app","target":"./.env","decision":"approved_once"}
```
