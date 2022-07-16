package inlog

type FmtPrinter struct {
}

func (f FmtPrinter) Debugf(format string, args ...any) { print("[DEBUG]", format, args...) }

func (f FmtPrinter) Infof(format string, args ...any) { print("[INFO]", format, args...) }

func (f FmtPrinter) Warnf(format string, args ...any) { print("[WARN]", format, args...) }

func (f FmtPrinter) Errorf(format string, args ...any) { print("[ERROR]", format, args...) }

func (f FmtPrinter) Panicf(format string, args ...any) { print("[PANIC]", format, args...) }

func (f FmtPrinter) Debug(args ...any) { print("[DEBUG]", "", args...) }

func (f FmtPrinter) Info(args ...any) { print("[INFO]", "", args...) }

func (f FmtPrinter) Warn(args ...any) { print("[WARN]", "", args...) }

func (f FmtPrinter) Error(args ...any) { print("[ERROR]", "", args...) }

func (f FmtPrinter) Panic(args ...any) { print("[PANIC]", "", args...) }

var _ IMinimLog = &FmtPrinter{}
