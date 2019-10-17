package scraper

import (
	"github.com/lruggieri/gonema/pkg/scraper/rarbg"
	"github.com/lruggieri/gonema/pkg/torrent"
)

type TorrentScraper interface{
	GetTorrentLinks(iResourceName,iResourceImdbID string) (oTorrents []torrent.Torrent, oErr error)
}

func New(iScraperType string) TorrentScraper{
	var chosenScraper TorrentScraper
	switch iScraperType{
	case rarbg.Name:{
		chosenScraper = &rarbg.Scraper{}
	}
	default:
		return nil
	}

	return chosenScraper
}