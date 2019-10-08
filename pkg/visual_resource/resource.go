package visual_resource

import (
	"encoding/json"
	"errors"
	"gitlab.com/ruggieri/gonema/pkg/scraper"
	"gitlab.com/ruggieri/gonema/pkg/scraper/rarbg"
	"gitlab.com/ruggieri/gonema/pkg/torrent"
	"time"
)

type Resource interface{
	Json() string


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
func GetMockResource() *mockResource{
	return &mockResource{}
}

type resourceImages struct {
	Small string `json:"small"`
	Big   string `json:"big"`
}
type resource struct {
	ImdbID            string            `json:"imdb_id"`
	Images            resourceImages    `json:"images"`
	ImdbScore         float64           `json:"imdb_score"` //out of 10
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
func(r *resource) Json() string{
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

type mockResource struct{
	resource
}
func (mr *mockResource) SetImdbID(iImdbID string) *mockResource{
	mr.ImdbID = iImdbID

	return mr
}
func (mr *mockResource) SetImages(iSmallImage, iBigImage string) *mockResource{
	mr.Images.Small = iSmallImage
	mr.Images.Big = iBigImage

	return mr
}
func (mr *mockResource) SetImdbScore(iImdbScore float64) *mockResource       { mr.ImdbScore = iImdbScore; return mr }
func (mr *mockResource) SetName(iName string) *mockResource                  { mr.Name = iName; return mr }
func (mr *mockResource) SetYear(iYear int) *mockResource                     { mr.Year = iYear; return mr }
func (mr *mockResource) SetReleaseDate(iReleaseDate time.Time) *mockResource { mr.ReleaseDate = iReleaseDate; return mr }
func (mr *mockResource) SetCategories(iCategories []string) *mockResource    { mr.Categories = iCategories; return mr }
func (mr *mockResource) SetPlot(iPlot string) *mockResource                  { mr.Plot = iPlot; return mr }
func (mr *mockResource) SetStars(iStars []string) *mockResource              { mr.Stars = iStars; return mr }
func (mr *mockResource) SetWriters(iWriters []string) *mockResource          { mr.Writers = iWriters; return mr }
func (mr *mockResource) SetDirectors(iDirectors []string) *mockResource      { mr.Directors = iDirectors; return mr }
func (mr *mockResource) SetAvailableTorrents(iAvailableTorrents []torrent.Torrent) *mockResource{
	mr.AvailableTorrents = iAvailableTorrents
	return mr
}
func (mr *mockResource) SetAvailableTorrent(iAvailableTorrent torrent.Torrent) *mockResource{
	mr.AvailableTorrents = append(mr.AvailableTorrents, iAvailableTorrent)
	return mr
}