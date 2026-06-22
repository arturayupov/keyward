# Installing keyward

keyward is a single static binary with no runtime dependencies.

## go install

```bash
go install github.com/arturayupov/keyward/cmd/keyward@latest
```

Ensure `$(go env GOPATH)/bin` is on your `PATH`.

## Pre-built binaries

Download the archive for your OS/arch from [Releases](https://github.com/arturayupov/keyward/releases), extract, and move `keyward` onto your `PATH`:

```bash
tar -xzf keyward_*_darwin_arm64.tar.gz
sudo mv keyward /usr/local/bin/
```

## Homebrew (macOS / Linux)

```bash
brew install arturayupov/tap/keyward
```

The Homebrew cask installs a signed-free binary and clears the macOS quarantine
attribute automatically, so it runs without a Gatekeeper prompt.

## Scoop (Windows)

```bash
scoop bucket add arturayupov https://github.com/arturayupov/scoop-bucket
scoop install keyward
```

## Build from source

```bash
git clone https://github.com/arturayupov/keyward
cd keyward
make build      # → bin/keyward
make install    # → $GOPATH/bin/keyward
```

Requires Go 1.22+.

## Per-OS notes

### macOS
The vault master key is stored in your **login Keychain**. The first time keyward
reads the key, macOS may show a one-time "keyward wants to use the keychain"
prompt — choose **Always Allow**. Release binaries are code-signed and notarized
so this prompt is minimal/absent. keyward resolves both the vault path
(`~/.keyward`) and the Keychain via `$HOME`; run it under your normal login (do
not override `$HOME`).

### Windows
The master key is stored in **Windows Credential Manager**. No extra setup.

### Linux
keyward needs a **Secret Service** provider (libsecret) — install and unlock one:

```bash
# Debian/Ubuntu
sudo apt install gnome-keyring libsecret-1-0
# Fedora
sudo dnf install gnome-keyring libsecret
```

For headless servers without a Secret Service, run inside a `dbus-run-session`
with `gnome-keyring-daemon` unlocked. (A passphrase-derived fallback is on the
roadmap.)

Native approval dialogs use `zenity` if present; otherwise keyward falls back to
a terminal prompt.

## Verify

```bash
keyward --help
keyward init      # creates ~/.keyward/vault.age and the keystore identity
```
