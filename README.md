# go-email-normalizer

> Fork of [dimuska139/go-email-normalizer](https://github.com/dimuska139/go-email-normalizer)

This is Golang library for providing a canonical representation of email address. It allows
to prevent multiple signups. `go-email-normalizer` contains some popular providers but you can easily append others.

> **RFC compliance note:** All normalization transformations (dot removal, `+` tag
> stripping, domain alias resolution) are applied only for specific, known providers
> where the behavior is documented. For unknown domains the local part is returned
> unchanged, consistent with RFC 5321 which treats the local part as opaque to
> external systems. See [NORMALIZATION.md](NORMALIZATION.md) for the full audit
> and per-provider rule reference.

## Usage

### Normalize

`Normalize` returns the canonical form of an email address as a plain string.

```go
package main

import (
	"fmt"
	"strings"
	normalizer "github.com/bobadilla-tech/go-email-normalizer"
)

type customRule struct{}

func (rule *customRule) ProcessUsername(username string) string {
	return strings.Replace(username, "-", "", -1)
}

func (rule *customRule) ProcessDomain(domain string) string {
	return domain
}

func main() {
	n := normalizer.NewNormalizer()
	fmt.Println(n.Normalize("vasya+pupkin@gmail.com")) // vasya@gmail.com
	fmt.Println(n.Normalize("t.e-St+vasya@gmail.com")) // te-st@gmail.com
	fmt.Println(n.Normalize("John+Brown@yahoo.com"))   // john+brown@yahoo.com
	fmt.Println(n.Normalize("John-Brown@yahoo.com"))   // john@yahoo.com
	fmt.Println(n.Normalize("t.e-St+@googlemail.com")) // te-st@gmail.com
	fmt.Println(n.Normalize("t.e-St+@google.com"))     // te-st@google.com

	n.AddRule("customrules.com", &customRule{})
	fmt.Println(n.Normalize(" tE-S-t@CustomRules.com.")) // tESt@customrules.com
}
```

### Normalize2

`Normalize2` accepts any string, validates it as an email address using an
RFC 5322-compatible regex (sourced from
[go-playground/validator](https://github.com/go-playground/validator)), and
returns a `NormalizeResult` paired with an `error`. The validator requires a
dot-separated domain (e.g. `gmail.com`), so inputs like `user@gmailcom` are
rejected. Trailing whitespace and trailing dots are stripped before validation,
so those are accepted. If the input is not a valid email address, the result is
zero-valued and the error is non-nil. When the call succeeds, the result pairs
the canonical address with a list of every transformation applied, in order.
Each `Change` value appears at most once.

```go
n := normalizer.NewNormalizer()

result, err := n.Normalize2("First.Last+tag@googlemail.com")
if err != nil {
	log.Fatal(err)
}
fmt.Println(result.Normalized) // firstlast@gmail.com
fmt.Println(result.Changes)
// [lowercase removed_dots removed_plus_tag canonicalized_domain]

result, err = n.Normalize2("User.Name_test-sub+spam@protonmail.com")
if err != nil {
	log.Fatal(err)
}
fmt.Println(result.Normalized) // usernametestsub@protonmail.com
fmt.Println(result.Changes)
// [lowercase removed_dots removed_underscores removed_hyphens removed_plus_tag]

result, err = n.Normalize2("Test+User-Name@ya.ru")
if err != nil {
	log.Fatal(err)
}
fmt.Println(result.Normalized) // testuser.name@yandex.ru
fmt.Println(result.Changes)
// [lowercase removed_plus_signs replaced_hyphens_with_dots canonicalized_domain]

// Invalid input — no "@" sign.
_, err = n.Normalize2("notanemail")
fmt.Println(err) // invalid email address: "notanemail"

// Invalid input — dotless domain is rejected.
_, err = n.Normalize2("not-an-email@gmailcom")
fmt.Println(err) // invalid email address: "not-an-email@gmailcom"
```

#### Change values

| Constant | Value | Produced by |
|---|---|---|
| `ChangeTrimmedWhitespace` | `trimmed_whitespace` | leading/trailing whitespace stripped |
| `ChangeRemovedTrailingDot` | `removed_trailing_dot` | trailing dot stripped from raw input |
| `ChangeLowercase` | `lowercase` | username or domain uppercased → lowercased |
| `ChangeRemovedDots` | `removed_dots` | dots removed from username (Google, Protonmail) |
| `ChangeRemovedUnderscores` | `removed_underscores` | underscores removed from username (Protonmail) |
| `ChangeRemovedHyphens` | `removed_hyphens` | hyphens removed from username (Protonmail) |
| `ChangeReplacedHyphensWithDots` | `replaced_hyphens_with_dots` | hyphens replaced with dots in username (Yandex) |
| `ChangeRemovedPlusTag` | `removed_plus_tag` | `+tag` subaddress stripped (Google, Apple, Fastmail, Protonmail) |
| `ChangeRemovedPlusSigns` | `removed_plus_signs` | all `+` characters removed (Microsoft, Rackspace, Rambler, Yandex, Zoho) |
| `ChangeRemovedSubaddress` | `removed_subaddress` | `-tag` subaddress stripped (Yahoo) |
| `ChangeCanonicalisedDomain` | `canonicalized_domain` | domain rewritten to canonical form (e.g. `googlemail.com` → `gmail.com`) |

## Supported providers

* Apple
* Fastmail
* Google
* Microsoft
* Protonmail
* Rackspace
* Rambler
* Yahoo
* Yandex
* Zoho

Also you can integrate other rules using `AddRule` function (see an example above)

For a detailed breakdown of which transformations each provider applies and the
RFC standards behind them, see [NORMALIZATION.md](NORMALIZATION.md).
