# keyward vs. alternatives

Honest comparisons. Every section names **when the alternative is the better
choice** — keyward is narrow on purpose, and trust comes from saying so.

The short version: tools like Keychain, `pass`, `sops`, and 1Password solve
**storing** secrets. keyward solves **handing a secret to an AI agent** — by
name, one at a time, only after you approve, without the model ever seeing the
value. It complements an existing store rather than replacing it.

---

## keyward vs. pasting keys into AI chat

The default today: you paste an API key into the prompt so the agent can use it.

| | Pasting into chat | keyward |
|---|---|---|
| Value enters model context | Yes | **No** |
| Lands in transcripts/logs | Yes | No |
| Per-use approval | No | Yes |
| Revocable without rotating | No (already exposed) | Yes (never exposed) |

**When pasting "wins":** a one-off throwaway key in a sandbox you'll delete. For
anything real, exposure is permanent — assume a pasted key is compromised.

---

## keyward vs. envchain

[envchain](https://github.com/sorah/envchain) stores environment variables in the
OS keychain and injects them into a process you launch.

**When envchain wins:** you just want to run a local process with secrets from the
Keychain and there's no AI agent involved. It's simple, battle-tested, and great
for `envchain stripe rails server`.

**When keyward wins:** an AI agent needs the key. envchain has no concept of a
request-by-name, no per-request approval, and no MCP interface — it injects
everything you grouped, to a process, with no human-in-the-loop gate. keyward
releases one named key per request, only after approval, and the agent never
receives the value.

---

## keyward vs. pass / sops / git-crypt

[pass](https://www.passwordstore.org/), [sops](https://github.com/getsops/sops),
and git-crypt are excellent at **encrypted storage and team sync** of secrets,
usually GPG/KMS-backed and git-friendly.

**When they win:** versioned, team-shared secret files; CI/CD pipelines; infra
secrets in a GitOps workflow. keyward is single-user and local; it does not try
to be your team's encrypted-config store.

**When keyward wins:** the consumer is an interactive AI agent, not a CI job. sops
decrypts a whole file into the environment; keyward brokers one key at a time with
a human approving each release, and keeps the value away from the model.

You can use both: keep your source of truth in sops/pass, `keyward import` what an
agent needs locally.

---

## keyward vs. 1Password CLI / Doppler / Infisical

These are mature secret managers (often cloud-backed) with CLIs, injection
(`op run`), sharing, and rich access policy.

**When they win:** you need a team vault, audit/compliance, SSO, rotation,
cross-machine sync, and broad integrations. keyward has none of that today.

**When keyward wins:** you specifically want an **agent-facing broker with
per-request human approval** where the value never reaches the model — a workflow
none of them target. keyward is local-first, free, open source, and MCP-native.

---

## Summary

keyward is not a replacement for your secret store. It's the **approval-gated
broker between your secrets and an AI agent** — the piece that's missing
everywhere else. If no AI agent is involved, you probably don't need it.
