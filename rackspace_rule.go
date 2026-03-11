package emailnormalizer

import "strings"

// RackspaceRule : email normalization rule for Rackspace
type RackspaceRule struct {
}

func (rule *RackspaceRule) ProcessUsername(username string) string {
	result := strings.ToLower(username)
	return strings.Replace(result, "+", "", -1)
}

func (rule *RackspaceRule) ProcessDomain(domain string) string {
	return domain
}

func (rule *RackspaceRule) ProcessUsernameWithChanges(username string) (string, []Change) {
	var changes []Change

	result := strings.ToLower(username)
	if result != username {
		changes = append(changes, ChangeLowercase)
	}

	withoutPlus := strings.Replace(result, "+", "", -1)
	if withoutPlus != result {
		changes = append(changes, ChangeRemovedPlusSigns)
	}

	return withoutPlus, changes
}

func (rule *RackspaceRule) ProcessDomainWithChanges(domain string) (string, []Change) {
	return rule.ProcessDomain(domain), nil
}
