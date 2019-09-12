package rarbg

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	"gonema/utils"
	"testing"
)

func TestGetTorrentLinks(t *testing.T) {
	utils.DebugActive = true
	utils.Logger.Level = logrus.DebugLevel
	finalTorrents,err := GetTorrentLinks("","tt6146584")
	if err != nil{
		t.Error(err)
		t.FailNow()
	}

	spew.Dump(finalTorrents)
}