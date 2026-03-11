# go-email-normalizer - email normalization for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/bobadilla-tech/go-email-normalizer.svg)](https://pkg.go.dev/github.com/bobadilla-tech/go-email-normalizer)
[![codecov](https://codecov.io/github/bobadilla-tech/go-email-normalizer/graph/badge.svg?token=E1TagDCBXc)](https://codecov.io/gh/bobadilla-tech/go-email-normalizer)
[![Go Report Card](https://goreportcard.com/badge/github.com/bobadilla-tech/go-email-normalizer)](https://goreportcard.com/report/github.com/bobadilla-tech/go-email-normalizer)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/bobadilla-tech/go-email-normalizer/blob/master/LICENSE)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)  

This is Golang library for providing a canonical representation of email address. It allows
to prevent multiple signups. `go-email-normalizer` contains some popular providers but you can easily append others.

## Download

```shell
go get -u github.com/bobadilla-tech/go-email-normalizer/v5
```

## Usage

### Normalize

`Normalize` returns the canonical form of an email address as a plain string.

```go
package main

import (
	"fmt"
	"strings"
	normalizer "github.com/bobadilla-tech/go-email-normalizer/v5"
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

`Normalize2` returns a `NormalizeResult` that pairs the canonical address with a
list of every transformation applied, in order. Each `Change` value appears at
most once.

```go
n := normalizer.NewNormalizer()

result := n.Normalize2("First.Last+tag@googlemail.com")
fmt.Println(result.Normalized) // firstlast@gmail.com
fmt.Println(result.Changes)
// [lowercase removed_dots removed_plus_tag canonicalized_domain]

result = n.Normalize2("User.Name_test-sub+spam@protonmail.com")
fmt.Println(result.Normalized) // usernametestsub@protonmail.com
fmt.Println(result.Changes)
// [lowercase removed_dots removed_underscores removed_hyphens removed_plus_tag]

result = n.Normalize2("Test+User-Name@ya.ru")
fmt.Println(result.Normalized) // testuser.name@yandex.ru
fmt.Println(result.Changes)
// [lowercase removed_plus_signs replaced_hyphens_with_dots canonicalized_domain]
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
