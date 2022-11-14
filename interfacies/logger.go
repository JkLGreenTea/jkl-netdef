package interfacies

// Logger - главный логгер.
type Logger interface {
	DEBUG(msg interface{})
	INFO(msg interface{})
	WARN(msg interface{})
	ERROR(msg interface{})
	FATAL(msg interface{})
}
