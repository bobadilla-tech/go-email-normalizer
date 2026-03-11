package emailnormalizer

import "strings"

// ZohoRule : email normalization rule for Zoho
type ZohoRule struct {
}

func (rule *ZohoRule) ProcessUsername(username string) string {
	result := strings.ToLower(username)
	return strings.Replace(result, "+", "", -1)
}

func (rule *ZohoRule) ProcessDomain(domain string) string {
	return domain
}

func (rule *ZohoRule) ProcessUsernameWithChanges(username string) (string, []Change) {
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

func (rule *ZohoRule) ProcessDomainWithChanges(domain string) (string, []Change) {
	return rule.ProcessDomain(domain), nil
}
