package interfacies

import (
	"JkLNetDef/engine/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HttpProxy interface {
	Listen() error
	Stop() error

	SetTitle(newTitle string)
	LoadTLSCertificate(domain, cert, key string) error

	GetID() primitive.ObjectID
	GetConfig() *config.ProxyConfig
	GetProtocol() string
	GetTitle() string
	GetType() string
}
