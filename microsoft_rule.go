package emailnormalizer

import "strings"

// MicrosoftRule : email normalization rule for Microsoft
type MicrosoftRule struct {
}

func (rule *MicrosoftRule) ProcessUsername(username string) string {
	result := strings.ToLower(username)
	return strings.Replace(result, "+", "", -1)
}

func (rule *MicrosoftRule) ProcessDomain(domain string) string {
	return domain
}

func (rule *MicrosoftRule) ProcessUsernameWithChanges(username string) (string, []Change) {
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

func (rule *MicrosoftRule) ProcessDomainWithChanges(domain string) (string, []Change) {
	return rule.ProcessDomain(domain), nil
}
