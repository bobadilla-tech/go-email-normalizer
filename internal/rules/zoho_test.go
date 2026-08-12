package rules

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestZohoUsername(t *testing.T) {
	rule := ZohoRule{}
	assert.Equal(t, "t.est", rule.ProcessUsername("t+.est"))
}

func TestZohoDomain(t *testing.T) {
	rule := ZohoRule{}
	assert.Equal(t, "zoho.com", rule.ProcessDomain("zoho.com"))
}

func TestZohoUsernameWithChanges(t *testing.T) {
	rule := ZohoRule{}
	result, changes := rule.ProcessUsernameWithChanges("Test+User")
	assert.Equal(t, "testuser", result)
	assert.Equal(t, []Change{ChangeLowercase, ChangeRemovedPlusSigns}, changes)

	result, changes = rule.ProcessUsernameWithChanges("testuser")
	assert.Equal(t, "testuser", result)
	assert.Nil(t, changes)
}

func TestZohoDomainWithChanges(t *testing.T) {
	rule := ZohoRule{}
	result, changes := rule.ProcessDomainWithChanges("zoho.com")
	assert.Equal(t, "zoho.com", result)
	assert.Nil(t, changes)
}
