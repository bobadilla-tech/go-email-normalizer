package emailnormalizer

import "strings"

// GoogleRule : email normalization rule for Google
type GoogleRule struct {
}

func (rule *GoogleRule) ProcessUsername(username string) string {
	result := strings.ToLower(username)
	result = strings.Replace(result, ".", "", -1)

	plusSignIndex := strings.Index(result, "+")
	if plusSignIndex != -1 {
		result = result[0:plusSignIndex]
	}

	return result
}

func (rule *GoogleRule) ProcessDomain(domain string) string {
	switch domain {
	case "google.com":
		return domain
	default:
		return "gmail.com" // googlemail.com/gmail.com => gmail.com
	}
}

func (rule *GoogleRule) ProcessUsernameWithChanges(username string) (string, []Change) {
	var changes []Change

	result := strings.ToLower(username)
	if result != username {
		changes = append(changes, ChangeLowercase)
	}

	withoutDots := strings.Replace(result, ".", "", -1)
	if withoutDots != result {
		changes = append(changes, ChangeRemovedDots)
	}
	result = withoutDots

	plusIndex := strings.Index(result, "+")
	if plusIndex != -1 {
		changes = append(changes, ChangeRemovedPlusTag)
		result = result[:plusIndex]
	}

	return result, changes
}

func (rule *GoogleRule) ProcessDomainWithChanges(domain string) (string, []Change) {
	result := rule.ProcessDomain(domain)
	if result != domain {
		return result, []Change{ChangeCanonicalisedDomain}
	}
	return result, nil
}
