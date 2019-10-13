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
}
func GetResources(iName, iImdbID string) (oResource Resource, oError error){

	resultingResources := resources{
		Resources:make([]*resource,0),
	}
	if len(iImdbID) > 0{
		resultingResources.Resources = append(resultingResources.Resources, &resource{
			ImdbID:iImdbID,
		})

		err := resultingResources.setInfo()
		if err != nil{
			return nil, err
		}
	}else{
		//TODO
		//search for resources with this name, or similar, and get slice of corresponding ImdbID
		//then for each ImdbID create a resource and set info
	}

	return &resultingResources, nil
}
func GetMockResource() *mockResource{
	return &mockResource{}
}

type resourceImages struct {
	Small string `json:"small"`
	Big   string `json:"big"`
}
type resources struct{
	Resources []*resource `json:"resources"`
}
func (rs *resources) Json() string{
	jsonResource, err := json.Marshal(rs)
	if err != nil{
		return ""
	}
	return string(jsonResource)
}
func (rs *resources) setInfo() error{
	for _,r := range rs.Resources{
		err := r.setInfo()
		if err != nil{
			return err
		}
	}
	return nil
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
	jsonResource, err := json.Marshal(r)
	if err != nil{
		return ""
	}
	return string(jsonResource)
}
func(r *resource) setInfo() error{

	//TODO fetch real information
	r.Images = resourceImages{Big:"https://cdn.suwalls.com/wallpapers/movies/marvin-15227-1920x1200.jpg"}

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