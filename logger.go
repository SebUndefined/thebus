package thebus

type Logger interface {
	Debug(msg string, kv ...any)
	Info(msg string, kv ...any)
	Warn(msg string, kv ...any)
	Error(msg string, kv ...any)
}

type noopLogger struct{}

func NoopLogger() Logger {
	return &noopLogger{}
}
func (*noopLogger) Info(msg string, kv ...any)   {}
func (*noopLogger) Error(msg string, kv ...any)  {}
func (*noopLogger) Debug(msg string, kv ...any)  {}
func (l *noopLogger) Warn(msg string, kv ...any) {}
