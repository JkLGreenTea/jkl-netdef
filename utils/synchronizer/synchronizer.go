package synchronizer

import (
	"JkLNetDef/engine/utils/waitgroup"
)

// Synchronizer - синхронизатор.
type Synchronizer struct {
	Logger     *Logger
	Proxy      *Proxy
	HttpServer *HttpServer
}

// Logger - глобальный синхронизатор.
type Logger struct {
	WaitGroup *waitgroup.WaitGroup
}

// Proxy - синхронизатор прокси.
type Proxy struct {
	WaitGroup *waitgroup.WaitGroup
}

// HttpServer - синхронизатор Http Api.
type HttpServer struct {
	WaitGroup  *waitgroup.WaitGroup
	RequestsWg *waitgroup.WaitGroup
}

// New - создание синхронизатора.
func New() *Synchronizer {
	return &Synchronizer{
		Proxy: &Proxy{
			WaitGroup: waitgroup.New(),
		},
		Logger: &Logger{
			WaitGroup: waitgroup.New(),
		},
		HttpServer: &HttpServer{
			WaitGroup:  waitgroup.New(),
			RequestsWg: waitgroup.New(),
		},
	}
}
