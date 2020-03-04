package injector

import (
	"github.com/sirupsen/logrus"
	"github.com/lruggieri/gonema/pkg/util"
	"testing"
)

func TestInjector_Run(t *testing.T) {
	inj := NewInjector("https://gonemapi.ruggieri.tech/resourceInfo",[]string{
		"imdbID=tt6146586",
	},100)


	util.DebugActive = false
	util.Logger.Level = logrus.DebugLevel
	inj.Run()
}