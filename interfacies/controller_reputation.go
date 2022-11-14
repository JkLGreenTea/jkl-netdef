package interfacies

import (
	"JkLNetDef/engine/models/controller_reputation"
	"net/http"
)

// ControllerReputation - контроллер репутации клиентов.
type ControllerReputation interface {
	Analise()
	StopAnalise()
	GetClientReputation(host, userAgent string) (*controller_reputation.Client, error)
	CheckClientReputation(host, userAgent string) (bool, error)
	ResetClientReputation(host, userAgent string)
	CheckHostReputation(host string) (bool, error)
	AddRequests(req *http.Request)
}
