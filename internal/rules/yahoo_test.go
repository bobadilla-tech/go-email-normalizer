package rules

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestYahooUsername(t *testing.T) {
	rule := YahooRule{}
	assert.Equal(t, "johnbrown", rule.ProcessUsername("JohnBrown"))
	assert.Equal(t, "john+brown", rule.ProcessUsername("John+Brown"))
	assert.Equal(t, "john", rule.ProcessUsername("John-Brown"))
}

func TestYahooDomain(t *testing.T) {
	rule := YahooRule{}
	assert.Equal(t, "yahoodns.net", rule.ProcessDomain("yahoodns.net"))
}

func TestYahooUsernameWithChanges(t *testing.T) {
	rule := YahooRule{}
	result, changes := rule.ProcessUsernameWithChanges("TestUser-subaddress")
	assert.Equal(t, "testuser", result)
	assert.Equal(t, []Change{ChangeLowercase, ChangeRemovedSubaddress}, changes)

	result, changes = rule.ProcessUsernameWithChanges("testuser")
	assert.Equal(t, "testuser", result)
	assert.Nil(t, changes)
}

func TestYahooDomainWithChanges(t *testing.T) {
	rule := YahooRule{}
	result, changes := rule.ProcessDomainWithChanges("yahoo.com")
	assert.Equal(t, "yahoo.com", result)
	assert.Nil(t, changes)
}
