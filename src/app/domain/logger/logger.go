package logger

type Logger interface {
	Info(args ...any)
	Infow(msg string, args ...any)
	Warn(args ...any)
	Error(args ...any)
}
