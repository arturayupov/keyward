# Troubleshooting & FAQ

## Keychain / keystore

### macOS asks "keyward wants to use the keychain"
Choose **Always Allow**. keyward stores the vault's master key in your login
Keychain; macOS prompts the first time an app accesses it. Unsigned builds (e.g.
`go install`) prompt more often — official **signed/notarized release binaries**
minimize this.

### "Keychain cannot be found" / keystore errors on macOS
macOS resolves the login Keychain via `$HOME`. Running keyward under an
overridden `$HOME` (sandbox, some CI setups) breaks keystore access. Run it under
your **normal login**, not a custom `$HOME`.

### Linux: "could not reach the OS keystore"
keyward needs a Secret Service (libsecret) provider, unlocked in your session:
```bash
sudo apt install gnome-keyring libsecret-1-0   # Debian/Ubuntu
sudo dnf install gnome-keyring libsecret       # Fedora
```
Headless server? Run inside an unlocked session:
```bash
dbus-run-session -- sh -c 'echo "" | gnome-keyring-daemon --unlock; keyward ...'
```
(A passphrase-derived fallback for keystore-less environments is on the roadmap.)

## Approval

### The approval dialog didn't appear
- **Linux:** native dialogs use `zenity`. If it's not installed, keyward falls
  back to a **terminal prompt** — make sure you're running in an interactive
  terminal, or `sudo apt install zenity`.
- **Over MCP with no GUI:** with no native dialog available, approval falls back
  to the terminal where `keyward serve-mcp` runs. Keep that terminal visible.

### A request was denied that I wanted to allow
Any dialog error, cancel, or timeout is treated as **Deny** (fail-closed) by
design. Just re-run the request and approve it.

## MCP integration

### The AI tool doesn't see keyward
1. Confirm `keyward` is on your `PATH` (`which keyward`).
2. Check the MCP config exactly:
   ```json
   { "mcpServers": { "keyward": { "command": "keyward", "args": ["serve-mcp"] } } }
   ```
   Claude Code: `~/.claude.json` or project `.mcp.json`. Cursor/Windsurf/Cline use
   the same `command`/`args`.
3. Fully restart the AI tool after editing config.
4. Sanity check the server runs: `keyward serve-mcp` should start and wait on stdio.

### "secret not found"
The name **and** namespace must match. List what's stored:
```bash
keyward ls          # all
keyward ls --ns shop-app
```
Namespaces come from the folder you imported from.

## Secrets & files

### Where does keyward store things?
| Path | What |
|---|---|
| `~/.keyward/vault.age` | encrypted secrets (`0600`) |
| OS keystore | the vault master key |
| `~/.keyward/audit.jsonl` | decision log (no values) |

### Do multi-line keys (PEM, service-account JSON) work?
Yes. Values with newlines, quotes, backslashes, and `#` are encoded so they
inject into a `.env` intact.

### Windows: injected file shows `-rw-rw-rw-`
Windows doesn't honor Unix `0600` mode bits; files inherit folder ACLs. Keep
injected `.env` files in a directory only your user can read. ACL-based hardening
is a roadmap item ([ROADMAP.md](ROADMAP.md)).

## Reset / uninstall

```bash
rm -rf ~/.keyward                                              # vault + audit
security delete-generic-password -s keyward -a vault-age-identity  # macOS keystore key
```
On Linux/Windows, remove the `keyward` entry via your Secret Service tool /
Credential Manager. Then `keyward init` starts fresh.

---

Still stuck? Open an issue with your OS, `keyward --version`, and the command
output (**never paste a secret value**).
