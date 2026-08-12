package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoogleUsername(t *testing.T) {
	rule := GoogleRule{}
	assert.Equal(t, "t", rule.ProcessUsername("t+.est"))
}

func TestGoogleDomain(t *testing.T) {
	rule := GoogleRule{}
	assert.Equal(t, "gmail.com", rule.ProcessDomain("googlemail.com"))
	assert.Equal(t, "gmail.com", rule.ProcessDomain("gmail.com"))
	assert.Equal(t, "google.com", rule.ProcessDomain("google.com"))
}

func TestGoogleUsernameWithChanges(t *testing.T) {
	rule := GoogleRule{}
	result, changes := rule.ProcessUsernameWithChanges("Test.User+tag")
	assert.Equal(t, "testuser", result)
	assert.Equal(t, []Change{ChangeLowercase, ChangeRemovedDots, ChangeRemovedPlusTag}, changes)

	result, changes = rule.ProcessUsernameWithChanges("testuser")
	assert.Equal(t, "testuser", result)
	assert.Nil(t, changes)
}

func TestGoogleDomainWithChanges(t *testing.T) {
	rule := GoogleRule{}
	result, changes := rule.ProcessDomainWithChanges("gmail.com")
	assert.Equal(t, "gmail.com", result)
	assert.Nil(t, changes)

	result, changes = rule.ProcessDomainWithChanges("googlemail.com")
	assert.Equal(t, "gmail.com", result)
	assert.Equal(t, []Change{ChangeCanonicalisedDomain}, changes)
}
