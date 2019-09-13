package visualResource

import (
	"encoding/json"
	"errors"
	"gonema/scraper"
	"gonema/scraper/rarbg"
	"gonema/torrent"
	"time"
)

type Resource interface{
	String() string


	setInfo() error
	setTorrentInfo() error
}
func GetResource(iName, iImdbID string) (oResource Resource, oError error){

	newResource :=  resource{
		ImdbID: iImdbID,
		Name:   iName,
	}

	err := newResource.setInfo()
	if err != nil{
		return nil,err
	}

	return &newResource, nil
}
type resource struct{
	ImdbID            string            `json:"imdb_id"`
	ImdbScore         float64           `json:"imdb_score"`//out of 10
	Name              string            `json:"title"`
	Year              int               `json:"year"`
	ReleaseDate       time.Time         `json:"release_date"`
	Categories        []string          `json:"categories"`
	Plot              string            `json:"plot"`
	Stars             []string          `json:"stars"`
	Writers           []string          `json:"writers"`
	Directors         []string          `json:"directors"`
	AvailableTorrents []torrent.Torrent `json:"available_torrents"`
}
func(r *resource) String() string{
	jsonMovie, err := json.Marshal(r)
	if err != nil{
		return ""
	}
	return string(jsonMovie)
}
func(r *resource) setInfo() error{

	err := r.setTorrentInfo()
	if err != nil{
		return err
	}

	return nil
}
func(r *resource) setTorrentInfo() error{
	scraperToUse := scraper.New(rarbg.Name)
	if scraperToUse == nil{
		return errors.New("invalid scraper was chosen")
	}

	torrentInfo,err := scraperToUse.GetTorrentLinks(r.Name,r.ImdbID)
	if err != nil{
		return err
	}

	r.AvailableTorrents = torrentInfo
	return nil
}