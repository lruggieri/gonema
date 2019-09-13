package scraper

import "gonema/torrent"

type TorrentScraper interface{
	GetTorrentLinks(iResourceName,iResourceImdbID string) (oTorrents []torrent.Torrent, oErr error)
}
