package rules

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProtonmailUsername(t *testing.T) {
	rule := ProtonmailRule{}
	assert.Equal(t, "t", rule.ProcessUsername("t+.est"))
}

func TestProtonmailDomain(t *testing.T) {
	rule := ProtonmailRule{}
	assert.Equal(t, "protonmail.ch", rule.ProcessDomain("protonmail.ch"))
}

func TestProtonmailUsernameWithChanges(t *testing.T) {
	rule := ProtonmailRule{}
	result, changes := rule.ProcessUsernameWithChanges("User.Name_test-sub+tag")
	assert.Equal(t, "usernametestsub", result)
	assert.Equal(t, []Change{
		ChangeLowercase,
		ChangeRemovedDots,
		ChangeRemovedUnderscores,
		ChangeRemovedHyphens,
		ChangeRemovedPlusTag,
	}, changes)

	result, changes = rule.ProcessUsernameWithChanges("username")
	assert.Equal(t, "username", result)
	assert.Nil(t, changes)
}

func TestProtonmailDomainWithChanges(t *testing.T) {
	rule := ProtonmailRule{}
	result, changes := rule.ProcessDomainWithChanges("protonmail.com")
	assert.Equal(t, "protonmail.com", result)
	assert.Nil(t, changes)
}

func TestProtonmailRule_ProcessUsername(t *testing.T) {
	type args struct {
		username string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "With plus sign",
			args: args{username: "t+est"},
			want: "t",
		}, {
			name: "With dot",
			args: args{username: "t.est"},
			want: "test",
		}, {
			name: "With underscore",
			args: args{username: "t_est"},
			want: "test",
		}, {
			name: "With underscore",
			args: args{username: "t-est"},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := &ProtonmailRule{}
			assert.Equal(t, tt.want, rule.ProcessUsername(tt.args.username))
		})
	}
}
