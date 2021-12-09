package goapp

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

func InitLogger() {
	logrus.SetLevel(LogLevel())
	logrus.SetReportCaller(true)

	formatter := &logrus.TextFormatter{
		TimestampFormat:        "02-01-2006 15:04:05",
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "-", fmt.Sprintf("%s:%d", formatFile(f.File), f.Line)
		},
	}

	logrus.SetFormatter(formatter)
	logrus.Info("Logger Initialized")
}

func formatFile(file string) string {
	path, err := os.Getwd()
	if err != nil {
		logrus.Fatal(err)
	}

	return strings.Replace(file, path, "~", 1)
}
