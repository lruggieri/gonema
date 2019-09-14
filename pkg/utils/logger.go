package utils

import (
	"github.com/sirupsen/logrus"
	"os"
)

/*
Checked in front of every DEBUG level logging, because the require A TON OF MEMORY due to the fact that they accept
a variadic of interfaces, which have to be translated in their basic types.
*/
var DebugActive bool

var Logger = logrus.New()

func init(){
	Logger.Formatter = new(logrus.JSONFormatter)
	Logger.Formatter = new(logrus.TextFormatter)
	Logger.Level = logrus.InfoLevel
	Logger.Out = os.Stdout
}