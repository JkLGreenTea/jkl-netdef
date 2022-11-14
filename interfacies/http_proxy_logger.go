package interfacies

import (
	"net/http"
	"time"
)

// HttpProxyLogger - proxy логгер http запросов.
type HttpProxyLogger interface {
	DEBUG(proxyTitle, clientLocation string, startTm, endTm time.Time, rw http.ResponseWriter, req *http.Request)
	INFO(proxyTitle, clientLocation string, startTm, endTm time.Time, rw http.ResponseWriter, req *http.Request)
	WARN(proxyTitle, clientLocation string, startTm, endTm time.Time, rw http.ResponseWriter, req *http.Request)
	ERROR(proxyTitle, clientLocation string, startTm, endTm time.Time, rw http.ResponseWriter, req *http.Request)
	FATAL(proxyTitle, clientLocation string, startTm, endTm time.Time, rw http.ResponseWriter, req *http.Request)

	FDEBUG(proxyTitle string, startTm, endTm time.Time, req *http.Request, res *http.Response)
	FINFO(proxyTitle string, startTm, endTm time.Time, req *http.Request, res *http.Response)
	FWARN(proxyTitle string, startTm, endTm time.Time, req *http.Request, res *http.Response)
	FERROR(proxyTitle string, startTm, endTm time.Time, req *http.Request, res *http.Response)
	FFATAL(proxyTitle string, startTm, endTm time.Time, req *http.Request, res *http.Response)
}
