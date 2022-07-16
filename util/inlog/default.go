package inlog

var defaultLogger IMinimLog = &FmtPrinter{}

// Debugf
// print Debugf logs
func Debugf(format string, args ...any) {
	defaultLogger.Debugf(format, args...)
}

// Infof
// print Infof logs
func Infof(format string, args ...any) {
	defaultLogger.Infof(format, args...)
}

// Warnf
// print Warnf logs
func Warnf(format string, args ...any) {
	defaultLogger.Warnf(format, args...)
}

// Errorf
// print Errorf logs
func Errorf(format string, args ...any) {
	defaultLogger.Errorf(format, args...)
}

// Panicf
// print Panicf logs
func Panicf(format string, args ...any) {
	defaultLogger.Panicf(format, args...)
}

// Debug
// print Debug logs
func Debug(args ...any) {
	defaultLogger.Debug(args...)
}

// Info
// print Info logs
func Info(args ...any) {
	defaultLogger.Info(args...)
}

// Warn
// print Warn logs
func Warn(args ...any) {
	defaultLogger.Warn(args...)
}

// Error
// print Error logs
func Error(args ...any) {
	defaultLogger.Error(args...)
}

// Panic
// print Panic logs
func Panic(args ...any) {
	defaultLogger.Panic(args...)
}
