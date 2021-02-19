package util

import (
	"github.com/sirupsen/logrus"
	"os"
)

//Logger is the main log variable. Will be called for every logging purposes.
var Logger = logrus.New()

func init() {
	Logger.Formatter = new(logrus.JSONFormatter)
	Logger.Formatter = new(logrus.TextFormatter)
	Logger.Level = logrus.InfoLevel
	Logger.Out = os.Stdout
}
