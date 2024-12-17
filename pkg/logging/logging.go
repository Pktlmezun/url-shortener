package logging

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

func Init() *logrus.Logger {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	l := logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{})
	mw := io.MultiWriter(os.Stdout, logFile)
	l.SetOutput(mw)
	l.SetLevel(logrus.InfoLevel)
	return l
}
