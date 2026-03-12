package emailnormalizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNormalizer_Normalize2_NoChanges verifies that a well-formed, already-canonical
// address produces an empty change list.
func TestNormalizer_Normalize2_NoChanges(t *testing.T) {
	n := NewNormalizer()
	result, err := n.Normalize2("test@gmail.com")
	assert.NoError(t, err)
	assert.Equal(t, "test@gmail.com", result.Normalized)
	assert.Empty(t, result.Changes)
}

func TestNormalizer_Normalize2_TrimmedWhitespace(t *testing.T) {
	n := NewNormalizer()
	result, err := n.Normalize2("  test@gmail.com  ")
	assert.NoError(t, err)
	assert.Equal(t, "test@gmail.com", result.Normalized)
	assert.Equal(t, []Change{ChangeTrimmedWhitespace}, result.Changes)
}

func TestNormalizer_Normalize2_RemovedTrailingDot(t *testing.T) {
	n := NewNormalizer()
	result, err := n.Normalize2("test@gmail.com.")
	assert.NoError(t, err)
	assert.Equal(t, "test@gmail.com", result.Normalized)
	assert.Equal(t, []Change{ChangeRemovedTrailingDot}, result.Changes)
}

// TestNormalizer_Normalize2_InvalidEmail verifies that a string without "@" is
// rejected with a non-nil error and a zero-valued NormalizeResult.
func TestNormalizer_Normalize2_InvalidEmail(t *testing.T) {
	n := NewNormalizer()
	result, err := n.Normalize2("notanemail")
	assert.Error(t, err)
	assert.Empty(t, result.Normalized)
	assert.Empty(t, result.Changes)
}

// TestNormalizer_Normalize2_InvalidEmailMultipleAt verifies that a string with
// multiple "@" characters is rejected with a non-nil error.
func TestNormalizer_Normalize2_InvalidEmailMultipleAt(t *testing.T) {
	n := NewNormalizer()
	result, err := n.Normalize2("a@b@gmail.com")
	assert.Error(t, err)
	assert.Empty(t, result.Normalized)
	assert.Empty(t, result.Changes)
}

// TestNormalizer_Normalize2_InvalidEmailEmptyLocal verifies that a string with an
// empty local part (e.g. "@gmail.com") is rejected with a non-nil error.
func TestNormalizer_Normalize2_InvalidEmailEmptyLocal(t *testing.T) {
	n := NewNormalizer()
	result, err := n.Normalize2("@gmail.com")
	assert.Error(t, err)
	assert.Empty(t, result.Normalized)
	assert.Empty(t, result.Changes)
}

// TestNormalizer_Normalize2_InvalidEmailEmptyDomain verifies that a string with an
// empty domain part (e.g. "user@") is rejected with a non-nil error.
func TestNormalizer_Normalize2_InvalidEmailEmptyDomain(t *testing.T) {
	n := NewNormalizer()
	result, err := n.Normalize2("user@")
	assert.Error(t, err)
	assert.Empty(t, result.Normalized)
	assert.Empty(t, result.Changes)
}

// TestNormalizer_Normalize2_InvalidEmailDotlessDomain verifies that a domain
// without a dot (e.g. "gmailcom") is rejected, closing the gap left by
// net/mail.ParseAddress which accepted such dotless domains as syntactically
// valid.
func TestNormalizer_Normalize2_InvalidEmailDotlessDomain(t *testing.T) {
	n := NewNormalizer()
	result, err := n.Normalize2("not-an-email@gmailcom")
	assert.Error(t, err)
	assert.Empty(t, result.Normalized)
	assert.Empty(t, result.Changes)
}

// TestNormalizer_Normalize2_UnknownDomain verifies that for an unknown domain only
// domain lowercasing is reported (username is left untouched per RFC 5321).
func TestNormalizer_Normalize2_UnknownDomain(t *testing.T) {
	n := NewNormalizer()
	result, err := n.Normalize2("User@EXAMPLE.COM")
	assert.NoError(t, err)
	assert.Equal(t, "User@example.com", result.Normalized)
	assert.Equal(t, []Change{ChangeLowercase}, result.Changes)
}

// TestNormalizer_Normalize2_UnknownDomainAlreadyLower verifies that no changes are
// reported when an unknown domain is already lowercase.
func TestNormalizer_Normalize2_UnknownDomainAlreadyLower(t *testing.T) {
	n := NewNormalizer()
	result, err := n.Normalize2("User@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "User@example.com", result.Normalized)
	assert.Empty(t, result.Changes)
}

// --- Google ---

func TestNormalizer_Normalize2_Google(t *testing.T) {
	n := NewNormalizer()
	result, err := n.Normalize2("Test.User+tag@gmail.com")
	assert.NoError(t, err)
	assert.Equal(t, "testuser@gmail.com", result.Normalized)
	assert.Equal(t, []Change{ChangeLowercase, ChangeRemovedDots, ChangeRemovedPlusTag}, result.Changes)
}

func TestNormalizer_Normalize2_GoogleNoDots(t *testing.T) {
	n := NewNormalizer()
	result, err := n.Normalize2("TestUser+tag@gmail.com")
	assert.NoError(t, err)
	assert.Equal(t, "testuser@gmail.com", result.Normalized)
	assert.Equal(t, []Change{ChangeLowercase, ChangeRemovedPlusTag}, result.Changes)
}

func TestNormalizer_Normalize2_GoogleCanonicalDomain(t *testing.T) {
	n := NewNormalizer()
	result, err := n.Normalize2("test@googlemail.com")
	assert.NoError(t, err)
	assert.Equal(t, "test@gmail.com", result.Normalized)
	assert.Equal(t, []Change{ChangeCanonicalisedDomain}, result.Changes)
}

// google.com is kept as-is (workspace accounts).
func TestNormalizer_Normalize2_GoogleDotCom(t *testing.T) {
	n := NewNormalizer()
	result, err := n.Normalize2("test@google.com")
	assert.NoError(t, err)
	assert.Equal(t, "test@google.com", result.Normalized)
	assert.Empty(t, result.Changes)
}

// --- Yahoo ---

func TestNormalizer_Normalize2_Yahoo(t *testing.T) {
	n := NewNormalizer()
	result, err := n.Normalize2("TestUser-subaddress@yahoo.com")
	assert.NoError(t, err)
	assert.Equal(t, "testuser@yahoo.com", result.Normalized)
	assert.Equal(t, []Change{ChangeLowercase, ChangeRemovedSubaddress}, result.Changes)
}

func TestNormalizer_Normalize2_YahooNoDash(t *testing.T) {
	n := NewNormalizer()
	result, err := n.Normalize2("TestUser@yahoo.com")
	assert.NoError(t, err)
	assert.Equal(t, "testuser@yahoo.com", result.Normalized)
	assert.Equal(t, []Change{ChangeLowercase}, result.Changes)
}

// --- Microsoft ---

func TestNormalizer_Normalize2_Microsoft(t *testing.T) {
	n := NewNormalizer()
	result, err := n.Normalize2("Test+User@hotmail.com")
	assert.NoError(t, err)
	assert.Equal(t, "testuser@hotmail.com", result.Normalized)
	assert.Equal(t, []Change{ChangeLowercase, ChangeRemovedPlusSigns}, result.Changes)
}

// --- Apple ---

func TestNormalizer_Normalize2_Apple(t *testing.T) {
	n := NewNormalizer()
	// me.com → icloud.com (canonicalized domain) + plus tag stripped
	result, err := n.Normalize2("TestUser+tag@me.com")
	assert.NoError(t, err)
	assert.Equal(t, "testuser@icloud.com", result.Normalized)
	assert.Equal(t, []Change{ChangeLowercase, ChangeRemovedPlusTag, ChangeCanonicalisedDomain}, result.Changes)
}

func TestNormalizer_Normalize2_AppleICloud(t *testing.T) {
	n := NewNormalizer()
	// icloud.com stays as icloud.com — no canonicalization change
	result, err := n.Normalize2("TestUser+tag@icloud.com")
	assert.NoError(t, err)
	assert.Equal(t, "testuser@icloud.com", result.Normalized)
	assert.Equal(t, []Change{ChangeLowercase, ChangeRemovedPlusTag}, result.Changes)
}

// --- Fastmail ---

func TestNormalizer_Normalize2_Fastmail(t *testing.T) {
	n := NewNormalizer()
	result, err := n.Normalize2("TestUser+tag@fastmail.com")
	assert.NoError(t, err)
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
	result, err := n.Normalize2("User.Name_test-sub+tag@protonmail.com")
	assert.NoError(t, err)
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
	result, err := n.Normalize2("testuser@protonmail.com")
	assert.NoError(t, err)
	assert.Equal(t, "testuser@protonmail.com", result.Normalized)
	assert.Empty(t, result.Changes)
}

// --- Rackspace ---

func TestNormalizer_Normalize2_Rackspace(t *testing.T) {
	n := NewNormalizer()
	result, err := n.Normalize2("Test+User@emailsrvr.com")
	assert.NoError(t, err)
	assert.Equal(t, "testuser@emailsrvr.com", result.Normalized)
	assert.Equal(t, []Change{ChangeLowercase, ChangeRemovedPlusSigns}, result.Changes)
}

// --- Rambler ---

func TestNormalizer_Normalize2_Rambler(t *testing.T) {
	n := NewNormalizer()
	result, err := n.Normalize2("Test+User@rambler.ru")
	assert.NoError(t, err)
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
	result, err := n.Normalize2("Test+User-Name@ya.ru")
	assert.NoError(t, err)
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
	result, err := n.Normalize2("testuser@yandex.ru")
	assert.NoError(t, err)
	assert.Equal(t, "testuser@yandex.ru", result.Normalized)
	assert.Empty(t, result.Changes)
}

// --- Zoho ---

func TestNormalizer_Normalize2_Zoho(t *testing.T) {
	n := NewNormalizer()
	result, err := n.Normalize2("Test+User@zoho.com")
	assert.NoError(t, err)
	assert.Equal(t, "testuser@zoho.com", result.Normalized)
	assert.Equal(t, []Change{ChangeLowercase, ChangeRemovedPlusSigns}, result.Changes)
}

// --- Deduplication ---

// TestNormalizer_Normalize2_LowercaseDedup verifies that ChangeLowercase is
// reported exactly once even when both the domain and username are uppercased.
func TestNormalizer_Normalize2_LowercaseDedup(t *testing.T) {
	n := NewNormalizer()
	result, err := n.Normalize2("TESTUSER@GMAIL.COM")
	assert.NoError(t, err)
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
// same normalized address as Normalize for a representative set of valid inputs.
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
		"User@EXAMPLE.COM",
		"  test@gmail.com  ",
		"test@gmail.com.",
		"test@googlemail.com",
		"test@google.com",
	}
	for _, email := range emails {
		result, err := n.Normalize2(email)
		assert.NoError(t, err, "expected no error for input: %q", email)
		assert.Equal(t, n.Normalize(email), result.Normalized,
			"Normalize2 and Normalize must agree for input: %q", email)
	}
}
