package log

func Debug(args ...interface{}) {
	log.Sugar().Debug(args...)
}

func Error(args ...interface{}) {
	log.Sugar().Error(args...)
}

func Info(args ...interface{}) {
	log.Sugar().Info(args...)
}

func Warn(args ...interface{}) {
	log.Sugar().Warn(args...)
}

func Fatal(args ...interface{}) {
	log.Sugar().Fatal(args...)
}

func Panic(args ...interface{}) {
	log.Sugar().Panic(args...)
}

func Debugf(template string, args ...interface{}) {
	log.Sugar().Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	log.Sugar().Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	log.Sugar().Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	log.Sugar().Errorf(template, args...)
}
func Panicf(template string, args ...interface{}) {
	log.Sugar().Panicf(template, args...)
}

func Fatalf(msg string, keysAndValues ...interface{}) {
	log.Sugar().Fatalf(msg, keysAndValues...)
}

func Debugln(args ...interface{}) {
	log.Sugar().Debugln(args...)
}

func Infoln(args ...interface{}) {
	log.Sugar().Infoln(args...)
}

func Warnln(args ...interface{}) {
	log.Sugar().Warnln(args...)
}

func Errorln(args ...interface{}) {
	log.Sugar().Errorln(args...)
}

func Panicln(args ...interface{}) {
	log.Sugar().Panicln(args...)
}

func Fatalln(args ...interface{}) {
	log.Sugar().Fatalln(args...)
}
