# Security

## Threat model

keyward exists to keep secret **values** away from AI agents while still letting
them use those secrets. The agent is treated as **untrusted with respect to
secret values** but trusted to run in your environment.

### Guarantees (enforced by automated tests)

1. **Value never reaches the model.** `request_key` injects the value into a
   target file and returns only `{status, name, target}`. `list_keys` returns
   names and namespaces only.
2. **Value never logged.** `~/.keyward/audit.jsonl` records tool, key, namespace,
   target, and decision — never the value. `ls`/`inject` never print it.
3. **Approval is out-of-band.** The approval dialog is fired by the keyward
   process via the OS, not rendered into the agent's token stream, so the agent
   cannot programmatically approve it.
4. **Fail closed.** Every approval backend returns **Deny** on error, cancel, or
   timeout. A denied request injects nothing.
5. **Encrypted at rest.** The vault is `age`-encrypted; the master key lives in
   the OS keystore, never in a plaintext file.

### What keyward does *not* defend against

- A compromised target file or project the agent can already read after
  injection (injection puts the value where your code needs it — by design).
- A malicious agent **socially engineering you** into approving a request. Read
  the dialog: it names the tool, key, namespace, and target path.
- Malware running as your user with keystore access. keyward raises the bar; it
  is not a substitute for OS-level security.

## Implementation notes

- **Keystore depends on `$HOME`.** On macOS the login Keychain is resolved via
  `$HOME`; running keyward under an overridden `$HOME` (e.g. a sandbox) can cause
  a "Keychain Not Found" error. Run under your normal login. The vault path
  (`~/.keyward`) follows the same `$HOME`.
- **Unsigned dev builds** may trigger a one-time OS-keystore access prompt.
  Release binaries are code-signed/notarized to minimize this.
- **Linux** requires a Secret Service (libsecret) provider; see
  [INSTALL.md](INSTALL.md).

## Reporting a vulnerability

Please **do not** open a public issue for security problems. Email the maintainer
(see the GitHub profile for `arturayupov`) with details and steps to reproduce.
You'll get an acknowledgement within a few days. Coordinated disclosure is
appreciated.
