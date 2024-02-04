package app

import (
	"github.com/sirupsen/logrus"
	"os"
)

func SetLogger() {

	logrus.SetLevel(logrus.DebugLevel)

	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logrus.SetOutput(os.Stdout)
}
