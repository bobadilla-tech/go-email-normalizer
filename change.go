package emailnormalizer

// Change is a string enum identifying a single transformation applied to an
// email address during normalization. It is returned by Normalize2.
type Change string

const (
	// ChangeTrimmedWhitespace is applied when leading or trailing whitespace is
	// removed from the raw input.
	ChangeTrimmedWhitespace Change = "trimmed_whitespace"

	// ChangeRemovedTrailingDot is applied when one or more trailing dots are
	// stripped from the raw input.
	ChangeRemovedTrailingDot Change = "removed_trailing_dot"

	// ChangeLowercase is applied when the username or domain contains uppercase
	// letters that are converted to lower case.
	ChangeLowercase Change = "lowercase"

	// ChangeRemovedDots is applied when dot characters are removed from the
	// username (e.g. Google, Protonmail).
	ChangeRemovedDots Change = "removed_dots"

	// ChangeRemovedUnderscores is applied when underscore characters are removed
	// from the username (e.g. Protonmail).
	ChangeRemovedUnderscores Change = "removed_underscores"

	// ChangeRemovedHyphens is applied when hyphen characters are removed from
	// the username (e.g. Protonmail).
	ChangeRemovedHyphens Change = "removed_hyphens"

	// ChangeReplacedHyphensWithDots is applied when hyphen characters in the
	// username are replaced with dots (e.g. Yandex).
	ChangeReplacedHyphensWithDots Change = "replaced_hyphens_with_dots"

	// ChangeRemovedPlusTag is applied when a plus-sign subaddress ("+tag") is
	// stripped from the end of the username (e.g. Google, Apple, Fastmail,
	// Protonmail).
	ChangeRemovedPlusTag Change = "removed_plus_tag"

	// ChangeRemovedPlusSigns is applied when plus-sign characters are removed
	// from the username without subaddress semantics — every "+" is deleted
	// regardless of position (e.g. Microsoft, Rackspace, Rambler, Yandex, Zoho).
	ChangeRemovedPlusSigns Change = "removed_plus_signs"

	// ChangeRemovedSubaddress is applied when a dash-delimited subaddress
	// ("-tag") is stripped from the end of the username (e.g. Yahoo).
	ChangeRemovedSubaddress Change = "removed_subaddress"

	// ChangeCanonicalisedDomain is applied when the domain is rewritten to its
	// canonical form (e.g. googlemail.com → gmail.com, me.com → icloud.com,
	// ya.ru → yandex.ru).
	ChangeCanonicalisedDomain Change = "canonicalized_domain"
)

// NormalizeResult is the return value of Normalize2.
type NormalizeResult struct {
	// Normalized is the canonical form of the email address.
	Normalized string

	// Changes lists every transformation that was applied, in the order they
	// were first detected. Each Change value appears at most once.
	Changes []Change
}
