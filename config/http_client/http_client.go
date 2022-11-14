package http_client

import "time"

// HTTPClient - HTTP клиент для проксирования запросов.
type HTTPClient struct {
	// Http client transport (http.Transport)
	Transport *Transport `json:"transport" bson:"transport" yaml:"transport" form:"transport" description:"Http client transport (http.Transport)"`

	// Http client dialer (net.Dialer)
	Dialer *Dialer `json:"dialer" bson:"dialer" yaml:"dialer" form:"dialer" description:"Http client dialer (net.Dialer)"`

	// Timeout задает ограничение по времени для запросов, сделанных клиентом.
	// Тайм-аут включает в себя время подключения, любые
	// перенаправления и чтение тела ответа.
	//
	// Тайм-аут, равный нулю, означает отсутствие тайм-аута
	// (по умолчанию 0).
	Timeout time.Duration `json:"timeout" bson:"timeout" yaml:"timeout" form:"timeout" description:"Timeout задает ограничение по времени для запросов, сделанных клиентом. Тайм-аут включает в себя время подключения, любые перенаправления и чтение тела ответа.Тайм-аут, равный нулю, означает отсутствие тайм-аута (по умолчанию 0)."`

	// Flush интервал (по умолчанию 0).
	FlushInterval time.Duration `json:"flush_interval" bson:"flush_interval" yaml:"flush_interval" form:"flush_interval" description:"Flush интервал (по умолчанию 0)."`

	// Версия протокола для входящих запросов сервера.
	//
	// Для клиентских запросов эти поля игнорируются.
	// Клиентский код всегда использует либо HTTP/1.1, либо HTTP/2
	// (по умолчанию 'HTTP/1.1').
	Proto string `json:"proto" bson:"proto" yaml:"proto" form:"proto" description:"Версия протокола для входящих запросов сервера.  Для клиентских запросов эти поля игнорируются. Клиентский код всегда использует либо HTTP/1.1, либо HTTP/2 (по умолчанию 'HTTP/1.1')."`

	// Главная версия (по умолчанию 1).
	ProtoMajor int `json:"proto_major" bson:"proto_major" yaml:"proto_major" form:"proto_major" description:"Главная версия (по умолчанию 1)."`

	// Второстепенная версия (по умолчанию 1).
	ProtoMinor int `json:"proto_minor" bson:"proto_minor" yaml:"proto_minor" form:"proto_minor" description:"Второстепенная версия (по умолчанию 1)."`
}
