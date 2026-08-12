package rules

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAppleUsername(t *testing.T) {
	rule := AppleRule{}
	assert.Equal(t, "johnbrown", rule.ProcessUsername("JohnBrown+test"))
}

func TestAppleDomain(t *testing.T) {
	rule := AppleRule{}
	assert.Equal(t, "icloud.com", rule.ProcessDomain("icloud.com"))
	assert.Equal(t, "icloud.com", rule.ProcessDomain("me.com"))
}

func TestAppleUsernameWithChanges(t *testing.T) {
	rule := AppleRule{}
	result, changes := rule.ProcessUsernameWithChanges("JohnBrown+test")
	assert.Equal(t, "johnbrown", result)
	assert.Equal(t, []Change{ChangeLowercase, ChangeRemovedPlusTag}, changes)

	result, changes = rule.ProcessUsernameWithChanges("johnbrown")
	assert.Equal(t, "johnbrown", result)
	assert.Nil(t, changes)
}

func TestAppleDomainWithChanges(t *testing.T) {
	rule := AppleRule{}
	result, changes := rule.ProcessDomainWithChanges("icloud.com")
	assert.Equal(t, "icloud.com", result)
	assert.Nil(t, changes)

	result, changes = rule.ProcessDomainWithChanges("me.com")
	assert.Equal(t, "icloud.com", result)
	assert.Equal(t, []Change{ChangeCanonicalisedDomain}, changes)
}
