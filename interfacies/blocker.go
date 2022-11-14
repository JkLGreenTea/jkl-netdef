package interfacies

import (
	blocker_model "JkLNetDef/engine/models/blocker"
	"JkLNetDef/engine/models/blocker/token"
)

// Blocker - блокировщик.
type Blocker interface {
	CheckHostInBlackList(host string) (bool, error)
	AddHostInBlackList(host, reason string) error
	RemoveHostFromBlackList(host string) error

	CheckHostInWhiteList(host string) (bool, error)
	AddHostInWhiteList(host, reason string) error
	RemoveHostFromWhiteList(host string) error

	CheckHostInHardWhiteList(host string) (bool, error)
	AddHostInHardWhiteList(host, reason string) error
	RemoveHostFromHardWhiteList(host string) error

	CheckHostInBanList(host string) (bool, error)
	AddHostInBanList(host, reason string, blockingTime int64) error
	RemoveHostFromBanList(host string) error

	CheckLocationInBlackList(location string) (bool, error)
	AddLocationInBlackList(location, reason string) error
	RemoveLocationFromBlackList(location string) error

	CheckLocationInWhiteList(location string) (bool, error)
	AddLocationInWhiteList(location, reason string) error
	RemoveLocationFromWhiteList(location string) error

	CheckUserAgentInWhiteList(agent string) (bool, error)
	AddUserAgentInWhiteList(agent, reason string) error
	RemoveUserAgentFromWhiteList(agent_ string) error

	CheckUserAgentInBlackList(agent string) (bool, error)
	AddUserAgentInBlackList(agent, reason string) error
	RemoveUserAgentFromBlackList(agent_ string) error

	CheckClientInClientListOnCaptchaCheck(host, userAgent string) (bool, error)
	GetClientInClientListOnCaptchaCheck(host, userAgent string) (*blocker_model.ClientOnCaptchaCheck, error)
	AddClientInClientListOnCaptchaCheck(host, userAgent, requestedUrl string) error
	RemoveClientFromClientListOnCaptchaCheck(host, userAgent string) error

	Location(host string) string

	NewToken() (string, error)
	GetToken(data string) (*token.Token, bool, error)
	RemoveToken(data string) error
	AddToken(tok *token.Token) error
}
