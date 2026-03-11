package emailnormalizer

// NormalizingRule : interface for all email normalization rules
type NormalizingRule interface {
	ProcessUsername(string) string
	ProcessDomain(string) string
}

// NormalizingRuleWithChanges is an optional extension of NormalizingRule.
// Rules that implement it allow Normalize2 to report the individual
// transformations applied to the username and domain. Rules added via
// AddRule that only satisfy NormalizingRule are still supported by
// Normalize2 — they will not report per-rule changes, but global
// pre-processing changes (whitespace trimming, trailing-dot removal,
// domain lowercasing) are always tracked.
type NormalizingRuleWithChanges interface {
	NormalizingRule
	ProcessUsernameWithChanges(username string) (result string, changes []Change)
	ProcessDomainWithChanges(domain string) (result string, changes []Change)
}
