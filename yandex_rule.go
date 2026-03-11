package emailnormalizer

import "strings"

// YandexRule : email normalization rule for Yandex
type YandexRule struct {
}

func (rule *YandexRule) ProcessUsername(username string) string {
	result := strings.ToLower(username)
	result = strings.Replace(result, "+", "", -1)
	return strings.Replace(result, "-", ".", -1)
}

func (rule *YandexRule) ProcessDomain(domain string) string {
	return "yandex.ru"
}

func (rule *YandexRule) ProcessUsernameWithChanges(username string) (string, []Change) {
	var changes []Change

	result := strings.ToLower(username)
	if result != username {
		changes = append(changes, ChangeLowercase)
	}

	withoutPlus := strings.Replace(result, "+", "", -1)
	if withoutPlus != result {
		changes = append(changes, ChangeRemovedPlusSigns)
	}
	result = withoutPlus

	withDots := strings.Replace(result, "-", ".", -1)
	if withDots != result {
		changes = append(changes, ChangeReplacedHyphensWithDots)
	}

	return withDots, changes
}

func (rule *YandexRule) ProcessDomainWithChanges(domain string) (string, []Change) {
	result := rule.ProcessDomain(domain)
	if result != domain {
		return result, []Change{ChangeCanonicalisedDomain}
	}
	return result, nil
}
