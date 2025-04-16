package ihplogger

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log = logrus.New()

func LogInit() {
	Log.SetOutput(&lumberjack.Logger{
		Filename:   "./log/system.log",
		MaxSize:    50, // MiB
		MaxBackups: 3,
		MaxAge:     30, // 30일 동안
		Compress:   true,
	})

	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
		ForceQuote:    true,
	})

	Log.SetLevel(logrus.DebugLevel)
	Log.SetReportCaller(true)
} 