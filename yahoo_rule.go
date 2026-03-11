package emailnormalizer

import "strings"

// YahooRule : email normalization rule for Yahoo
type YahooRule struct {
}

func (rule *YahooRule) ProcessUsername(username string) string {
	result := strings.ToLower(username)

	subaddressingIndex := strings.Index(result, "-")
	if subaddressingIndex != -1 {
		result = result[0:subaddressingIndex]
	}

	return result
}

func (rule *YahooRule) ProcessDomain(domain string) string {
	return domain
}

func (rule *YahooRule) ProcessUsernameWithChanges(username string) (string, []Change) {
	var changes []Change

	result := strings.ToLower(username)
	if result != username {
		changes = append(changes, ChangeLowercase)
	}

	dashIndex := strings.Index(result, "-")
	if dashIndex != -1 {
		changes = append(changes, ChangeRemovedSubaddress)
		result = result[:dashIndex]
	}

	return result, changes
}

func (rule *YahooRule) ProcessDomainWithChanges(domain string) (string, []Change) {
	return rule.ProcessDomain(domain), nil
}
