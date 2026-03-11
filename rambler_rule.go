package emailnormalizer

import "strings"

// RamblerRule : email normalization rule for Rambler
type RamblerRule struct {
}

func (rule *RamblerRule) ProcessUsername(username string) string {
	result := strings.ToLower(username)
	return strings.Replace(result, "+", "", -1)
}

func (rule *RamblerRule) ProcessDomain(domain string) string {
	return domain
}

func (rule *RamblerRule) ProcessUsernameWithChanges(username string) (string, []Change) {
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

func (rule *RamblerRule) ProcessDomainWithChanges(domain string) (string, []Change) {
	return rule.ProcessDomain(domain), nil
}
