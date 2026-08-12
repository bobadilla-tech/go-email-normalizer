package rules

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRamblerUsername(t *testing.T) {
	rule := RamblerRule{}
	assert.Equal(t, "t.est", rule.ProcessUsername("t+.est"))
}

func TestRamblerDomain(t *testing.T) {
	rule := RamblerRule{}
	assert.Equal(t, "lenta.ru", rule.ProcessDomain("lenta.ru"))
	assert.Equal(t, "rambler.ru", rule.ProcessDomain("rambler.ru"))
}

func TestRamblerUsernameWithChanges(t *testing.T) {
	rule := RamblerRule{}
	result, changes := rule.ProcessUsernameWithChanges("Test+User")
	assert.Equal(t, "testuser", result)
	assert.Equal(t, []Change{ChangeLowercase, ChangeRemovedPlusSigns}, changes)

	result, changes = rule.ProcessUsernameWithChanges("testuser")
	assert.Equal(t, "testuser", result)
	assert.Nil(t, changes)
}

func TestRamblerDomainWithChanges(t *testing.T) {
	rule := RamblerRule{}
	result, changes := rule.ProcessDomainWithChanges("rambler.ru")
	assert.Equal(t, "rambler.ru", result)
	assert.Nil(t, changes)
}
