package log

import (
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func InitLogger(logLevel logrus.Level) error {
	// Create new logger
	logger := logrus.New()
	logger.SetLevel(logLevel)

	// Add filesystem hook in order to log to files
	pathMap := lfshook.PathMap{
		logrus.InfoLevel:  "./application.log",
		logrus.WarnLevel:  "./application.log",
		logrus.ErrorLevel: "./error.log",
		logrus.FatalLevel: "./error.log",
	}
	var hook *lfshook.LfsHook = lfshook.NewHook(
		pathMap,
		&logrus.JSONFormatter{PrettyPrint: false},
	)
	logger.Hooks.Add(hook)
	log = logger
	return nil
}

func formMessage(opts ...string) string {
	// If only one opt is provided, return it
	if len(opts) == 1 {
		return opts[0]
	}
	output := ""
	for _, opt := range opts {
		output += opt
	}
	return output
}

func genericLogger(level logrus.Level, opts ...string) {
	// Will only log message if no opts are provided in order to avoid redundant variable allocation
	if len(opts) == 0 {
		return
	}
	log.Log(level, formMessage(opts...))
}

func Trace(opts ...string) {
	genericLogger(logrus.TraceLevel, opts...)
}

func Debug(opts ...string) {
	genericLogger(logrus.DebugLevel, opts...)
}

func Info(opts ...string) {
	genericLogger(logrus.InfoLevel, opts...)
}

func Warn(opts ...string) {
	genericLogger(logrus.WarnLevel, opts...)
}

func Error(opts ...string) {
	genericLogger(logrus.ErrorLevel, opts...)
}

func Fatal(opts ...string) {
	genericLogger(logrus.FatalLevel, opts...)
	output := ""
	for _, opt := range opts {
		output += opt
	}
	panic(output)
}
