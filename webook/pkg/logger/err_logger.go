package logger

type ErrLogger struct {
	ZapLogger LoggerV1
}

func NewErrLogger(l LoggerV1) ErrLogger {
	return ErrLogger{
		ZapLogger: l,
	}
}

func (e *ErrLogger) HandleError(err error, msg string, args ...Field) {
	if err != nil {
		e.ZapLogger.Error(msg, Field{
			Key: "origErr",
			Val: err,
		})
	}
}
