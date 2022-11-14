package engine

import (
	"JkLNetDef/engine/config"
	"JkLNetDef/engine/interfacies"
)

// State - состояние движка.
type State interface {
	run() error
	stop() error
	newProxy(cfg *config.ProxyConfig) (interfacies.HttpProxy, error)
}

// States - состояния движка.
type States struct {
	Run  *engineStateRun
	Stop *engineStateStop
}
