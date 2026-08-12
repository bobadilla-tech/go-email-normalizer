package emailnormalizer

import (
	"strings"

	"github.com/bobadilla-tech/go-email-normalizer/internal/rules"
)

// Normalizer : main library object for normalization emails
type Normalizer struct {
	rules map[string]NormalizingRule
}

// NewNormalizer : creates Normalizer instance
func NewNormalizer() *Normalizer {
	registry := make(map[string]NormalizingRule)

	microsoftRule := &rules.MicrosoftRule{}
	for _, domain := range rules.MicrosoftDomains {
		registry[domain] = microsoftRule
	}

	yahooRule := &rules.YahooRule{}
	for _, domain := range rules.YahooDomains {
		registry[domain] = yahooRule
	}

	googleRule := &rules.GoogleRule{}
	for _, domain := range rules.GoogleDomains {
		registry[domain] = googleRule
	}

	fastmailRule := &rules.FastmailRule{}
	for _, domain := range rules.FastmailDomains {
		registry[domain] = fastmailRule
	}

	ramblerRule := &rules.RamblerRule{}
	for _, domain := range rules.RamblerDomains {
		registry[domain] = ramblerRule
	}

	yandexRule := &rules.YandexRule{}
	for _, domain := range rules.YandexDomains {
		registry[domain] = yandexRule
	}

	protonmailRule := &rules.ProtonmailRule{}
	for _, domain := range rules.ProtonmailDomains {
		registry[domain] = protonmailRule
	}

	appleRule := &rules.AppleRule{}
	for _, domain := range rules.AppleDomains {
		registry[domain] = appleRule
	}

	registry["emailsrvr.com"] = &rules.RackspaceRule{}
	registry["zoho.com"] = &rules.ZohoRule{}
	return &Normalizer{rules: registry}
}

// AddRule : appends custom normalization rule
func (n *Normalizer) AddRule(domain string, strategy NormalizingRule) {
	n.rules[domain] = strategy
}

// Normalize : converts email to canonical form
func (n *Normalizer) Normalize(email string) string {
	prepared := strings.TrimSpace(email)
	prepared = strings.TrimRight(prepared, ".")

	parts := strings.Split(prepared, "@")
	if len(parts) != 2 {
		return prepared
	}

	username := parts[0]                // The first part of the address may be case sensitive (RFC 5336)
	domain := strings.ToLower(parts[1]) // Domain names are case-insensitive (RFC 4343)

	if rule, ok := n.rules[domain]; ok {
		return rule.ProcessUsername(username) + "@" + rule.ProcessDomain(domain)
	}

	return username + "@" + domain
}

// Normalize2 converts any input string to a canonical email address and
// returns a NormalizeResult that includes both the normalized address and
// every transformation applied.
//
// Unlike Normalize, Normalize2 validates the input using an RFC 5322-compatible
// email regex (sourced from github.com/go-playground/validator). It requires a
// dot-separated domain (e.g. "gmail.com"), so inputs like "user@gmailcom" are
// rejected. Trailing whitespace and trailing dots are stripped before
// validation, so those are handled correctly. If the input is not a valid email
// address, a zero-valued NormalizeResult and a non-nil error are returned.
//
// Rules registered via AddRule that only implement NormalizingRule (not
// NormalizingRuleWithChanges) are still fully supported — the normalized
// address will be correct, but per-rule changes will not be reported.
// Pre-processing changes (whitespace trimming, trailing-dot removal, domain
// lowercasing) are always tracked regardless of the rule type.
func (n *Normalizer) Normalize2(email string) (NormalizeResult, error) {
	seen := make(map[Change]bool)
	var changes []Change

	addChange := func(c Change) {
		if !seen[c] {
			seen[c] = true
			changes = append(changes, c)
		}
	}

	// Step 1: Trim leading/trailing whitespace.
	prepared := strings.TrimSpace(email)
	if prepared != email {
		addChange(ChangeTrimmedWhitespace)
	}

	// Step 2: Strip trailing dots.
	trimmed := strings.TrimRight(prepared, ".")
	if trimmed != prepared {
		addChange(ChangeRemovedTrailingDot)
	}
	prepared = trimmed

	// Step 3: Validate using the go-playground/validator email regex and extract the address.
	if err := ValidateEmail(prepared); err != nil {
		return NormalizeResult{}, err
	}
	parts := strings.SplitN(prepared, "@", 2)

	username := parts[0]                // The first part of the address may be case sensitive (RFC 5336)
	domain := strings.ToLower(parts[1]) // Domain names are case-insensitive (RFC 4343)
	if domain != parts[1] {
		addChange(ChangeLowercase)
	}

	if rule, ok := n.rules[domain]; ok {
		if detailed, ok := rule.(NormalizingRuleWithChanges); ok {
			normalizedUsername, usernameChanges := detailed.ProcessUsernameWithChanges(username)
			normalizedDomain, domainChanges := detailed.ProcessDomainWithChanges(domain)
			for _, c := range usernameChanges {
				addChange(c)
			}
			for _, c := range domainChanges {
				addChange(c)
			}
			return NormalizeResult{
				Normalized: normalizedUsername + "@" + normalizedDomain,
				Changes:    changes,
			}, nil
		}
		// Fallback: rule does not implement NormalizingRuleWithChanges.
		return NormalizeResult{
			Normalized: rule.ProcessUsername(username) + "@" + rule.ProcessDomain(domain),
			Changes:    changes,
		}, nil
	}

	return NormalizeResult{Normalized: username + "@" + domain, Changes: changes}, nil
}
