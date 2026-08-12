package rules

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFastmailUsername(t *testing.T) {
	rule := FastmailRule{}
	assert.Equal(t, "johnbrown", rule.ProcessUsername("JohnBrown+test"))
}

func TestFastmailDomain(t *testing.T) {
	rule := FastmailRule{}
	assert.Equal(t, "fastmail.com", rule.ProcessDomain("fastmail.com"))
	assert.Equal(t, "messagingengine.com", rule.ProcessDomain("messagingengine.com"))
}

func TestFastmailUsernameWithChanges(t *testing.T) {
	rule := FastmailRule{}
	result, changes := rule.ProcessUsernameWithChanges("JohnBrown+test")
	assert.Equal(t, "johnbrown", result)
	assert.Equal(t, []Change{ChangeLowercase, ChangeRemovedPlusTag}, changes)

	result, changes = rule.ProcessUsernameWithChanges("johnbrown")
	assert.Equal(t, "johnbrown", result)
	assert.Nil(t, changes)
}

func TestFastmailDomainWithChanges(t *testing.T) {
	rule := FastmailRule{}
	result, changes := rule.ProcessDomainWithChanges("fastmail.com")
	assert.Equal(t, "fastmail.com", result)
	assert.Nil(t, changes)
}
