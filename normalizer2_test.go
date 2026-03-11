package emailnormalizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNormalizer_Normalize2_NoChanges verifies that a well-formed, already-canonical
// address produces an empty change list.
func TestNormalizer_Normalize2_NoChanges(t *testing.T) {
	n := NewNormalizer()
	result := n.Normalize2("test@gmail.com")
	assert.Equal(t, "test@gmail.com", result.Normalized)
	assert.Empty(t, result.Changes)
}

func TestNormalizer_Normalize2_TrimmedWhitespace(t *testing.T) {
	n := NewNormalizer()
	result := n.Normalize2("  test@gmail.com  ")
	assert.Equal(t, "test@gmail.com", result.Normalized)
	assert.Equal(t, []Change{ChangeTrimmedWhitespace}, result.Changes)
}

func TestNormalizer_Normalize2_RemovedTrailingDot(t *testing.T) {
	n := NewNormalizer()
	result := n.Normalize2("test@gmail.com.")
	assert.Equal(t, "test@gmail.com", result.Normalized)
	assert.Equal(t, []Change{ChangeRemovedTrailingDot}, result.Changes)
}

func TestNormalizer_Normalize2_InvalidEmail(t *testing.T) {
	n := NewNormalizer()
	result := n.Normalize2("notanemail")
	assert.Equal(t, "notanemail", result.Normalized)
	assert.Empty(t, result.Changes)
}

// TestNormalizer_Normalize2_UnknownDomain verifies that for an unknown domain only
// domain lowercasing is reported (username is left untouched per RFC 5321).
func TestNormalizer_Normalize2_UnknownDomain(t *testing.T) {
	n := NewNormalizer()
	result := n.Normalize2("User@EXAMPLE.COM")
	assert.Equal(t, "User@example.com", result.Normalized)
	assert.Equal(t, []Change{ChangeLowercase}, result.Changes)
}

// TestNormalizer_Normalize2_UnknownDomainAlreadyLower verifies that no changes are
// reported when an unknown domain is already lowercase.
func TestNormalizer_Normalize2_UnknownDomainAlreadyLower(t *testing.T) {
	n := NewNormalizer()
	result := n.Normalize2("User@example.com")
	assert.Equal(t, "User@example.com", result.Normalized)
	assert.Empty(t, result.Changes)
}

// --- Google ---

func TestNormalizer_Normalize2_Google(t *testing.T) {
	n := NewNormalizer()
	result := n.Normalize2("Test.User+tag@gmail.com")
	assert.Equal(t, "testuser@gmail.com", result.Normalized)
	assert.Equal(t, []Change{ChangeLowercase, ChangeRemovedDots, ChangeRemovedPlusTag}, result.Changes)
}

func TestNormalizer_Normalize2_GoogleNoDots(t *testing.T) {
	n := NewNormalizer()
	result := n.Normalize2("TestUser+tag@gmail.com")
	assert.Equal(t, "testuser@gmail.com", result.Normalized)
	assert.Equal(t, []Change{ChangeLowercase, ChangeRemovedPlusTag}, result.Changes)
}

func TestNormalizer_Normalize2_GoogleCanonicalDomain(t *testing.T) {
	n := NewNormalizer()
	result := n.Normalize2("test@googlemail.com")
	assert.Equal(t, "test@gmail.com", result.Normalized)
	assert.Equal(t, []Change{ChangeCanonicalisedDomain}, result.Changes)
}

// google.com is kept as-is (workspace accounts).
func TestNormalizer_Normalize2_GoogleDotCom(t *testing.T) {
	n := NewNormalizer()
	result := n.Normalize2("test@google.com")
	assert.Equal(t, "test@google.com", result.Normalized)
	assert.Empty(t, result.Changes)
}

// --- Yahoo ---

func TestNormalizer_Normalize2_Yahoo(t *testing.T) {
	n := NewNormalizer()
	result := n.Normalize2("TestUser-subaddress@yahoo.com")
	assert.Equal(t, "testuser@yahoo.com", result.Normalized)
	assert.Equal(t, []Change{ChangeLowercase, ChangeRemovedSubaddress}, result.Changes)
}

func TestNormalizer_Normalize2_YahooNoDash(t *testing.T) {
	n := NewNormalizer()
	result := n.Normalize2("TestUser@yahoo.com")
	assert.Equal(t, "testuser@yahoo.com", result.Normalized)
	assert.Equal(t, []Change{ChangeLowercase}, result.Changes)
}

// --- Microsoft ---

func TestNormalizer_Normalize2_Microsoft(t *testing.T) {
	n := NewNormalizer()
	result := n.Normalize2("Test+User@hotmail.com")
	assert.Equal(t, "testuser@hotmail.com", result.Normalized)
	assert.Equal(t, []Change{ChangeLowercase, ChangeRemovedPlusSigns}, result.Changes)
}

// --- Apple ---

func TestNormalizer_Normalize2_Apple(t *testing.T) {
	n := NewNormalizer()
	// me.com → icloud.com (canonicalized domain) + plus tag stripped
	result := n.Normalize2("TestUser+tag@me.com")
	assert.Equal(t, "testuser@icloud.com", result.Normalized)
	assert.Equal(t, []Change{ChangeLowercase, ChangeRemovedPlusTag, ChangeCanonicalisedDomain}, result.Changes)
}

func TestNormalizer_Normalize2_AppleICloud(t *testing.T) {
	n := NewNormalizer()
	// icloud.com stays as icloud.com — no canonicalization change
	result := n.Normalize2("TestUser+tag@icloud.com")
	assert.Equal(t, "testuser@icloud.com", result.Normalized)
	assert.Equal(t, []Change{ChangeLowercase, ChangeRemovedPlusTag}, result.Changes)
}

// --- Fastmail ---

func TestNormalizer_Normalize2_Fastmail(t *testing.T) {
	n := NewNormalizer()
	result := n.Normalize2("TestUser+tag@fastmail.com")
	assert.Equal(t, "testuser@fastmail.com", result.Normalized)
	assert.Equal(t, []Change{ChangeLowercase, ChangeRemovedPlusTag}, result.Changes)
}

// --- Protonmail ---

// TestNormalizer_Normalize2_Protonmail exercises all five username changes at once.
func TestNormalizer_Normalize2_Protonmail(t *testing.T) {
	n := NewNormalizer()
	// lower:             "user.name_test-sub+tag"
	// remove ".":        "username_test-sub+tag"
	// remove "_":        "usernametest-sub+tag"
	// remove "-":        "usernametestsub+tag"
	// strip "+tag":      "usernametestsub"
	result := n.Normalize2("User.Name_test-sub+tag@protonmail.com")
	assert.Equal(t, "usernametestsub@protonmail.com", result.Normalized)
	assert.Equal(t, []Change{
		ChangeLowercase,
		ChangeRemovedDots,
		ChangeRemovedUnderscores,
		ChangeRemovedHyphens,
		ChangeRemovedPlusTag,
	}, result.Changes)
}

func TestNormalizer_Normalize2_ProtonmailNoPunctuation(t *testing.T) {
	n := NewNormalizer()
	result := n.Normalize2("testuser@protonmail.com")
	assert.Equal(t, "testuser@protonmail.com", result.Normalized)
	assert.Empty(t, result.Changes)
}

// --- Rackspace ---

func TestNormalizer_Normalize2_Rackspace(t *testing.T) {
	n := NewNormalizer()
	result := n.Normalize2("Test+User@emailsrvr.com")
	assert.Equal(t, "testuser@emailsrvr.com", result.Normalized)
	assert.Equal(t, []Change{ChangeLowercase, ChangeRemovedPlusSigns}, result.Changes)
}

// --- Rambler ---

func TestNormalizer_Normalize2_Rambler(t *testing.T) {
	n := NewNormalizer()
	result := n.Normalize2("Test+User@rambler.ru")
	assert.Equal(t, "testuser@rambler.ru", result.Normalized)
	assert.Equal(t, []Change{ChangeLowercase, ChangeRemovedPlusSigns}, result.Changes)
}

// --- Yandex ---

func TestNormalizer_Normalize2_Yandex(t *testing.T) {
	n := NewNormalizer()
	// lower:                "test+user-name"
	// remove "+":           "testuser-name"
	// replace "-" with ".": "testuser.name"
	// domain: ya.ru → yandex.ru
	result := n.Normalize2("Test+User-Name@ya.ru")
	assert.Equal(t, "testuser.name@yandex.ru", result.Normalized)
	assert.Equal(t, []Change{
		ChangeLowercase,
		ChangeRemovedPlusSigns,
		ChangeReplacedHyphensWithDots,
		ChangeCanonicalisedDomain,
	}, result.Changes)
}

func TestNormalizer_Normalize2_YandexPrimaryDomain(t *testing.T) {
	n := NewNormalizer()
	// yandex.ru → yandex.ru: no canonicalization change
	result := n.Normalize2("testuser@yandex.ru")
	assert.Equal(t, "testuser@yandex.ru", result.Normalized)
	assert.Empty(t, result.Changes)
}

// --- Zoho ---

func TestNormalizer_Normalize2_Zoho(t *testing.T) {
	n := NewNormalizer()
	result := n.Normalize2("Test+User@zoho.com")
	assert.Equal(t, "testuser@zoho.com", result.Normalized)
	assert.Equal(t, []Change{ChangeLowercase, ChangeRemovedPlusSigns}, result.Changes)
}

// --- Deduplication ---

// TestNormalizer_Normalize2_LowercaseDedup verifies that ChangeLowercase is
// reported exactly once even when both the domain and username are uppercased.
func TestNormalizer_Normalize2_LowercaseDedup(t *testing.T) {
	n := NewNormalizer()
	result := n.Normalize2("TESTUSER@GMAIL.COM")
	assert.Equal(t, "testuser@gmail.com", result.Normalized)
	count := 0
	for _, c := range result.Changes {
		if c == ChangeLowercase {
			count++
		}
	}
	assert.Equal(t, 1, count, "ChangeLowercase must appear exactly once")
}

// --- Parity ---

// TestNormalizer_Normalize2_Parity verifies that Normalize2 always produces the
// same normalized address as Normalize for a representative set of inputs.
func TestNormalizer_Normalize2_Parity(t *testing.T) {
	n := NewNormalizer()
	emails := []string{
		"Test.User+tag@gmail.com",
		"TestUser-subaddress@yahoo.com",
		"Test+User@hotmail.com",
		"TestUser+tag@me.com",
		"User.Name_test-sub+tag@protonmail.com",
		"Test+User-Name@ya.ru",
		"Test+User@zoho.com",
		"notanemail",
		"User@EXAMPLE.COM",
		"  test@gmail.com  ",
		"test@gmail.com.",
		"test@googlemail.com",
		"test@google.com",
	}
	for _, email := range emails {
		assert.Equal(t, n.Normalize(email), n.Normalize2(email).Normalized,
			"Normalize2 and Normalize must agree for input: %q", email)
	}
}
