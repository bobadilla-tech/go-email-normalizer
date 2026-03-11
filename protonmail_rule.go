package emailnormalizer

import "strings"

// ProtonmailRule : email normalization rule for Protonmail
type ProtonmailRule struct {
}

func (rule *ProtonmailRule) ProcessUsername(username string) string {
	result := strings.ToLower(username)

	charsToReplace := []string{
		".",
		"_",
		"-",
	}

	for _, char := range charsToReplace {
		result = strings.Replace(result, char, "", -1)
	}

	plusSignIndex := strings.Index(result, "+")
	if plusSignIndex != -1 {
		result = result[0:plusSignIndex]
	}

	return result
}

func (rule *ProtonmailRule) ProcessDomain(domain string) string {
	return domain
}

func (rule *ProtonmailRule) ProcessUsernameWithChanges(username string) (string, []Change) {
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

	withoutUnderscores := strings.Replace(result, "_", "", -1)
	if withoutUnderscores != result {
		changes = append(changes, ChangeRemovedUnderscores)
	}
	result = withoutUnderscores

	withoutHyphens := strings.Replace(result, "-", "", -1)
	if withoutHyphens != result {
		changes = append(changes, ChangeRemovedHyphens)
	}
	result = withoutHyphens

	plusIndex := strings.Index(result, "+")
	if plusIndex != -1 {
		changes = append(changes, ChangeRemovedPlusTag)
		result = result[:plusIndex]
	}

	return result, changes
}

func (rule *ProtonmailRule) ProcessDomainWithChanges(domain string) (string, []Change) {
	return rule.ProcessDomain(domain), nil
}
