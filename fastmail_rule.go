package emailnormalizer

import "strings"

// FastmailRule : email normalization rule for Fastmail
type FastmailRule struct {
}

func (rule *FastmailRule) ProcessUsername(username string) string {
	result := strings.ToLower(username)

	// Remove sub-addressing part (RFC 5233)
	plusSignIndex := strings.Index(result, "+")
	if plusSignIndex != -1 {
		result = result[0:plusSignIndex]
	}

	return result
}

func (rule *FastmailRule) ProcessDomain(domain string) string {
	return domain
}

func (rule *FastmailRule) ProcessUsernameWithChanges(username string) (string, []Change) {
	var changes []Change

	result := strings.ToLower(username)
	if result != username {
		changes = append(changes, ChangeLowercase)
	}

	plusIndex := strings.Index(result, "+")
	if plusIndex != -1 {
		changes = append(changes, ChangeRemovedPlusTag)
		result = result[:plusIndex]
	}

	return result, changes
}

func (rule *FastmailRule) ProcessDomainWithChanges(domain string) (string, []Change) {
	return rule.ProcessDomain(domain), nil
}
