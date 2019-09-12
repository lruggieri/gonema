package movie

import (
	"gonema/torrent"
	"time"
)

type Movie struct{
	ImdbID string
	ImdbScore float64 //out of 10
	Title string
	Year int
	ReleaseDate time.Time
	Duration time.Duration
	Categories []Category
	Plot string
	Stars []string
	Writers []string
	Directors []string
	AvailableTorrents []torrent.Torrent
}