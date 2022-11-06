package inlog

type FmtPrinter struct{}

func (f FmtPrinter) Debugf(format string, args ...any) { innerPrint("[DEBUG]", format, args...) }

func (f FmtPrinter) Infof(format string, args ...any) { innerPrint("[INFO]", format, args...) }

func (f FmtPrinter) Warnf(format string, args ...any) { innerPrint("[WARN]", format, args...) }

func (f FmtPrinter) Errorf(format string, args ...any) { innerPrint("[ERROR]", format, args...) }

func (f FmtPrinter) Panicf(format string, args ...any) { innerPrint("[PANIC]", format, args...) }

func (f FmtPrinter) Debug(args ...any) { innerPrint("[DEBUG]", "", args...) }

func (f FmtPrinter) Info(args ...any) { innerPrint("[INFO]", "", args...) }

func (f FmtPrinter) Warn(args ...any) { innerPrint("[WARN]", "", args...) }

func (f FmtPrinter) Error(args ...any) { innerPrint("[ERROR]", "", args...) }

func (f FmtPrinter) Panic(args ...any) { innerPrint("[PANIC]", "", args...) }

var _ IMinimLog = &FmtPrinter{}
