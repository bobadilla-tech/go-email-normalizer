package rules

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMicrosoftUsername(t *testing.T) {
	rule := MicrosoftRule{}
	assert.Equal(t, "t.est", rule.ProcessUsername("t+.est"))
}

func TestMicrosoftDomain(t *testing.T) {
	rule := MicrosoftRule{}
	assert.Equal(t, "microsoft.com", rule.ProcessDomain("microsoft.com"))
	assert.Equal(t, "live.com", rule.ProcessDomain("live.com"))
}

func TestMicrosoftUsernameWithChanges(t *testing.T) {
	rule := MicrosoftRule{}
	result, changes := rule.ProcessUsernameWithChanges("Test+User")
	assert.Equal(t, "testuser", result)
	assert.Equal(t, []Change{ChangeLowercase, ChangeRemovedPlusSigns}, changes)

	result, changes = rule.ProcessUsernameWithChanges("testuser")
	assert.Equal(t, "testuser", result)
	assert.Nil(t, changes)
}

func TestMicrosoftDomainWithChanges(t *testing.T) {
	rule := MicrosoftRule{}
	result, changes := rule.ProcessDomainWithChanges("hotmail.com")
	assert.Equal(t, "hotmail.com", result)
	assert.Nil(t, changes)
}
