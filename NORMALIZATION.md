# Email Normalization: Assumptions and Provider Rules

## RFC Background

Email normalization must respect several relevant RFCs:

- **RFC 5321** (SMTP) defines the local part of an email address as opaque to the
  transport layer. The receiving mail server is the sole authority on whether two
  local-part strings identify the same mailbox.
- **RFC 5322** (Internet Message Format) specifies the syntax of email addresses but
  does not prescribe any equivalence rules for the local part.
- **RFC 4343** defines that domain names are case-insensitive.
- **RFC 5233** defines the *Sieve subaddress extension* ("`+`" tagging), which is
  optional and must be enabled by the receiving server. It is not universally
  supported.

**Consequence:** No normalization transformation that modifies the local part is
universally correct. Every such transformation must be justified by documented,
provider-specific behavior.

---

## What This Library Does (and Does Not Do)

### Always applied (RFC-compliant)

| Transformation | Rationale |
|---|---|
| Strip leading/trailing whitespace | Syntactically invalid; safe to remove |
| Strip trailing dot from domain | Per DNS conventions; safe to remove |
| Lowercase the domain | RFC 4343: domain names are case-insensitive |

### Never applied globally

| Transformation | Why not |
|---|---|
| Remove dots from local part | Gmail-specific; other providers treat `john.doe` and `johndoe` as different accounts |
| Strip `+tag` subaddress | RFC 5233 is optional; not universally supported |
| Remove `+` signs entirely | Provider-specific (some servers reject `+` in the local part) |
| Remove or replace `_`, `-` | Provider-specific |
| Normalize domain aliases | Only explicit, documented mappings are used |

### Unknown / unregistered domains

For email addresses whose domain is not registered in the normalizer, the local
part is returned **unchanged**. Only the domain is lowercased (RFC 4343). This
follows the RFC 5321 principle that the local part is opaque to external systems.

```
user@example.com        → user@example.com     (no change to local part)
User@EXAMPLE.COM        → User@example.com     (only domain lowercased)
User+tag@example.com    → User+tag@example.com (+ tag NOT stripped)
john.doe@example.com    → john.doe@example.com (dot NOT removed)
```

---

## Per-Provider Normalization Rules

### Google (`gmail.com`, `googlemail.com`, `google.com`)

| Transformation | Example |
|---|---|
| Lowercase local part | `John` → `john` |
| Remove all dots from local part | `john.doe` → `johndoe` |
| Strip `+tag` subaddress | `john+spam` → `john` |
| Canonicalize domain | `googlemail.com` → `gmail.com` |
| `google.com` preserved as-is | `user@google.com` → `user@google.com` (Workspace accounts) |

**References:** [Gmail Help – About dots in addresses](https://support.google.com/mail/answer/7436150),
[Gmail Help – Create tasks from emails (plus addressing)](https://support.google.com/mail/answer/22370)

### Yahoo (`yahoo.com`, `ymail.com`, and regional domains)

| Transformation | Example |
|---|---|
| Lowercase local part | `John` → `john` |
| Strip `-tag` subaddress (dash, not plus) | `john-spam` → `john` |

**References:** [Yahoo Help – Disposable email addresses](https://help.yahoo.com/kb/SLN2028.html)

### Microsoft (`hotmail.com`, `outlook.com`, `live.com`, `msn.com`, and regional variants)

| Transformation | Example |
|---|---|
| Lowercase local part | `John` → `john` |
| Remove all `+` characters | `john+tag` → `johntag` |

Note: Microsoft's servers do not support `+` subaddressing; the `+` sign in the
local part is stripped server-side rather than being used as a subaddress separator.

### Apple (`icloud.com`, `me.com`, `mac.com`)

| Transformation | Example |
|---|---|
| Lowercase local part | `John` → `john` |
| Strip `+tag` subaddress | `john+spam` → `john` |
| Canonicalize domain | `me.com` → `icloud.com`, `mac.com` → `icloud.com` |

**References:** [Apple Support – iCloud Mail: Use plus addressing](https://support.apple.com/guide/icloud/use-plus-addressing-mm6b1a0955/icloud)
(Note: Apple's `+tag` support is also described in their iCloud Mail documentation; alias domains `me.com` and `mac.com` are legacy names for `icloud.com`.)

### Fastmail (`fastmail.com`, `fastmail.fm`, `messagingengine.com`)

| Transformation | Example |
|---|---|
| Lowercase local part | `John` → `john` |
| Strip `+tag` subaddress | `john+spam` → `john` |

**References:** [Fastmail Help – Plus addressing](https://www.fastmail.help/hc/en-us/articles/1500000280261)

### Protonmail (`protonmail.com`, `protonmail.ch`, `proton.me`, `pm.me`)

| Transformation | Example |
|---|---|
| Lowercase local part | `John` → `john` |
| Remove dots from local part | `john.doe` → `johndoe` |
| Remove underscores from local part | `john_doe` → `johndoe` |
| Remove hyphens from local part | `john-doe` → `johndoe` |
| Strip `+tag` subaddress | `john+spam` → `john` |

**References:** [Proton Help – What is my Proton Mail address](https://proton.me/support/email-alias-addresses)

### Yandex (`yandex.ru`, `ya.ru`, and international Yandex domains)

| Transformation | Example |
|---|---|
| Lowercase local part | `John` → `john` |
| Remove all `+` characters | `john+tag` → `johntag` |
| Replace hyphens with dots | `john-doe` → `john.doe` |
| Canonicalize domain | `ya.ru` → `yandex.ru` (and all other Yandex domains → `yandex.ru`) |

**References:** [Yandex Help – Email address rules](https://yandex.com/support/mail/mail-clients.html)

### Rackspace (`emailsrvr.com`)

| Transformation | Example |
|---|---|
| Lowercase local part | `John` → `john` |
| Remove all `+` characters | `john+tag` → `johntag` |

### Rambler (`rambler.ru`, `lenta.ru`, `autorambler.ru`, `myrambler.ru`, `ro.ru`)

| Transformation | Example |
|---|---|
| Lowercase local part | `John` → `john` |
| Remove all `+` characters | `john+tag` → `johntag` |

### Zoho (`zoho.com`)

| Transformation | Example |
|---|---|
| Lowercase local part | `John` → `john` |
| Remove all `+` characters | `john+tag` → `johntag` |

---

## Domain Alias Normalization

Domain alias resolution maps non-canonical domains to their canonical equivalents.
**Only explicit, provider-documented mappings are used.** No DNS lookups or
heuristics are performed.

| Input domain | Canonical domain | Provider |
|---|---|---|
| `googlemail.com` | `gmail.com` | Google |
| `me.com` | `icloud.com` | Apple |
| `mac.com` | `icloud.com` | Apple |
| `ya.ru` | `yandex.ru` | Yandex |
| All other `yandex.*` | `yandex.ru` | Yandex |

---

## Answers to the RFC Audit Questions

| Question | Answer |
|---|---|
| Are dots removed globally? | **No.** Only for Google and Protonmail. |
| Are `+` tags stripped for all domains? | **No.** `+tag` stripping applies only to Google, Apple, Fastmail, and Protonmail. For Microsoft, Rackspace, Rambler, Yandex, and Zoho, all `+` characters are removed (because those providers do not support `+` in local parts at all). Unknown domains are untouched. |
| Is normalization provider-specific? | **Yes.** Every transformation is scoped to a registered provider. |
| Is the canonical domain used in the output? | **Yes.** `Normalize` and `Normalize2` always build the output address from `ProcessDomain(domain)`, which returns the canonical domain. |
| Are both the input and canonical domain exposed? | `Normalize` returns a plain string; the caller can compare it with the input to detect domain changes. `Normalize2` returns a `NormalizeResult` whose `Changes` field includes `ChangeCanonicalisedDomain` when the domain was rewritten, allowing callers to detect this transformation. |

---

## Adding a New Provider

1. Add the provider's domains to `domains.go` as `var <provider>Domains = []string{...}`.
2. Create `<provider>_rule.go` implementing `NormalizingRule` (and optionally
   `NormalizingRuleWithChanges` for `Normalize2` change tracking).
3. Register every known domain for the provider in `NewNormalizer()` in `normalizer.go`.
4. Create `<provider>_rule_test.go` covering `ProcessUsername` and `ProcessDomain`.
5. Document the normalization behavior in the table above and cite the provider's
   official documentation.
