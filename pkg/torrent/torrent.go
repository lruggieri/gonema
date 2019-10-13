package torrent

import (
	"encoding/json"
	"github.com/nmmh/magneturi/magneturi"
)

type Torrent struct{
	MagnetLink magneturi.MagnetURI `json:"magnet_link"`
	Quality string `json:"quality"`
	Length string `json:"length"`
	Resolution string `json:"resolution"`
	Size string `json:"size"`
	Sound string `json:"sound"`
	Codec string `json:"codec"`
	Name string `json:"name"`
	Seeders int `json:"seeders"`
	Leechers int `json:"leechers"`
	Subtitles string `json:"subtitles"`
}
//solution to customize how some Torrent's field is marshaled
func (t *Torrent) MarshalJSON() ([]byte, error){
	type alias Torrent
	return json.Marshal(&struct{
		MagnetLink string `json:"magnet_link"`
		*alias
	}{
		MagnetLink:t.MagnetLink.String(),
		alias:(*alias)(t),
	})
}