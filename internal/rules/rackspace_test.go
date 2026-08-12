package rules

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRackspaceUsername(t *testing.T) {
	rule := RackspaceRule{}
	assert.Equal(t, "t.est", rule.ProcessUsername("t+.est"))
}

func TestRackspaceDomain(t *testing.T) {
	rule := RackspaceRule{}
	assert.Equal(t, "emailsrvr.com", rule.ProcessDomain("emailsrvr.com"))
}

func TestRackspaceUsernameWithChanges(t *testing.T) {
	rule := RackspaceRule{}
	result, changes := rule.ProcessUsernameWithChanges("Test+User")
	assert.Equal(t, "testuser", result)
	assert.Equal(t, []Change{ChangeLowercase, ChangeRemovedPlusSigns}, changes)

	result, changes = rule.ProcessUsernameWithChanges("testuser")
	assert.Equal(t, "testuser", result)
	assert.Nil(t, changes)
}

func TestRackspaceDomainWithChanges(t *testing.T) {
	rule := RackspaceRule{}
	result, changes := rule.ProcessDomainWithChanges("emailsrvr.com")
	assert.Equal(t, "emailsrvr.com", result)
	assert.Nil(t, changes)
}
