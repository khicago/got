package inlog

var defaultLogger IMinimLog = &FmtPrinter{}

// Debugf
// innerPrint Debugf logs
func Debugf(format string, args ...any) {
	defaultLogger.Debugf(format, args...)
}

// Infof
// innerPrint Infof logs
func Infof(format string, args ...any) {
	defaultLogger.Infof(format, args...)
}

// Warnf
// innerPrint Warnf logs
func Warnf(format string, args ...any) {
	defaultLogger.Warnf(format, args...)
}

// Errorf
// innerPrint Errorf logs
func Errorf(format string, args ...any) {
	defaultLogger.Errorf(format, args...)
}

// Panicf
// innerPrint Panicf logs
func Panicf(format string, args ...any) {
	defaultLogger.Panicf(format, args...)
}

// Debug
// innerPrint Debug logs
func Debug(args ...any) {
	defaultLogger.Debug(args...)
}

// Info
// innerPrint Info logs
func Info(args ...any) {
	defaultLogger.Info(args...)
}

// Warn
// innerPrint Warn logs
func Warn(args ...any) {
	defaultLogger.Warn(args...)
}

// Error
// innerPrint Error logs
func Error(args ...any) {
	defaultLogger.Error(args...)
}

// Panic
// innerPrint Panic logs
func Panic(args ...any) {
	defaultLogger.Panic(args...)
}
