package emailnormalizer

import "strings"

// AppleRule : email normalization rule for Apple
type AppleRule struct {
}

func (rule *AppleRule) ProcessUsername(username string) string {
	result := strings.ToLower(username)
	// Remove sub-addressing part (RFC 5233)
	plusSignIndex := strings.Index(result, "+")
	if plusSignIndex != -1 {
		result = result[0:plusSignIndex]
	}

	return result
}

func (rule *AppleRule) ProcessDomain(domain string) string {
	return "icloud.com"
}

func (rule *AppleRule) ProcessUsernameWithChanges(username string) (string, []Change) {
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

func (rule *AppleRule) ProcessDomainWithChanges(domain string) (string, []Change) {
	result := rule.ProcessDomain(domain)
	if result != domain {
		return result, []Change{ChangeCanonicalisedDomain}
	}
	return result, nil
}
