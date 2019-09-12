package torrent

import "github.com/nmmh/magneturi/magneturi"

type Torrent struct{
	MagnetLink magneturi.MagnetURI
	Quality string
	Resolution string
	Sound string
	Codec string
	Name string
	Seeders int
	Leechers int
}