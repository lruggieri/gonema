package rarbg

import (
	"github.com/sirupsen/logrus"
	"testing"
	"gonema/utils"
)

func TestGetTorrentLinks(t *testing.T) {
	utils.DebugActive = true
	utils.Logger.Level = logrus.DebugLevel
	err := GetTorrentLinks("","tt6146586")
	if err != nil{
		t.Error(err)
		t.FailNow()
	}
}