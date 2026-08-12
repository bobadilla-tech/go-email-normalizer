package rules

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestYandexUsername(t *testing.T) {
	rule := YandexRule{}
	assert.Equal(t, "t.est", rule.ProcessUsername("t+.est"))
	assert.Equal(t, "t.est", rule.ProcessUsername("t-est"))
}

func TestYandexDomain(t *testing.T) {
	rule := YandexRule{}
	assert.Equal(t, "yandex.ru", rule.ProcessDomain("yandex.ru"))
	assert.Equal(t, "yandex.ru", rule.ProcessDomain("ya.ru"))
	assert.Equal(t, "yandex.ru", rule.ProcessDomain("narod.ru"))
	assert.Equal(t, "yandex.ru", rule.ProcessDomain("yandex.com"))
	assert.Equal(t, "yandex.ru", rule.ProcessDomain("yandex.by"))
	assert.Equal(t, "yandex.ru", rule.ProcessDomain("yandex.ua"))
	assert.Equal(t, "yandex.ru", rule.ProcessDomain("yandex.kz"))
}

func TestYandexUsernameWithChanges(t *testing.T) {
	rule := YandexRule{}
	result, changes := rule.ProcessUsernameWithChanges("Test+User-Name")
	assert.Equal(t, "testuser.name", result)
	assert.Equal(t, []Change{ChangeLowercase, ChangeRemovedPlusSigns, ChangeReplacedHyphensWithDots}, changes)

	result, changes = rule.ProcessUsernameWithChanges("testuser")
	assert.Equal(t, "testuser", result)
	assert.Nil(t, changes)
}

func TestYandexDomainWithChanges(t *testing.T) {
	rule := YandexRule{}
	result, changes := rule.ProcessDomainWithChanges("yandex.ru")
	assert.Equal(t, "yandex.ru", result)
	assert.Nil(t, changes)

	result, changes = rule.ProcessDomainWithChanges("ya.ru")
	assert.Equal(t, "yandex.ru", result)
	assert.Equal(t, []Change{ChangeCanonicalisedDomain}, changes)
}
