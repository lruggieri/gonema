package controller

import (
	"github.com/nmmh/magneturi/magneturi"
	"gitlab.com/ruggieri/gonema/pkg/torrent"
	"gitlab.com/ruggieri/gonema/pkg/visual_resource"
)

func GetResourceInfo(resourceName, resourceImdbID string) ([]visual_resource.Resource, error) {
	//TODO really
	mockResource := visual_resource.GetMockResource().SetImdbID("tt6146586").SetImages(
		"https://m.media-amazon.com/images/M/MV5BMDg2YzI0ODctYjliMy00NTU0LTkxODYtYTNkNjQwMzVmOTcxXkEyXkFqcGdeQXVyNjg2NjQwMDQ@._V1_UX182_CR0,0,182,268_AL_.jpg",
		"https://m.media-amazon.com/images/M/MV5BMDg2YzI0ODctYjliMy00NTU0LTkxODYtYTNkNjQwMzVmOTcxXkEyXkFqcGdeQXVyNjg2NjQwMDQ@._V1_SY1000_CR0,0,648,1000_AL_.jpg",
	).SetName("John Wick 3 - Parabellum").SetYear(2019).SetCategories([]string{"Badass", "Immortal", "Kratos with nice hair and guns, lots of guns"},
	).SetAvailableTorrents([]torrent.Torrent{
		{
			MagnetLink: magneturi.MagnetURI{},
			Quality:    "SuperDuper",
			Resolution: "42K",
			Sound:      "H264",
			Name:       "Keanu4Evah",
			Seeders:    42,
			Leechers:   42,
		},
		{
			MagnetLink: magneturi.MagnetURI{},
			Quality:    "SuperDuper2",
			Resolution: "42K",
			Sound:      "H264",
			Name:       "Keanu4Evah2",
			Seeders:    42,
			Leechers:   42,
		},
		{
			MagnetLink: magneturi.MagnetURI{},
			Quality:    "SuperDuper3",
			Resolution: "42K",
			Sound:      "H264",
			Name:       "Keanu4Evah3",
			Seeders:    42,
			Leechers:   42,
		},
		{
			MagnetLink: magneturi.MagnetURI{},
			Quality:    "SuperDuper4",
			Resolution: "42K",
			Sound:      "H264",
			Name:       "Keanu4Evah4",
			Seeders:    42,
			Leechers:   42,
		},
		{
			MagnetLink: magneturi.MagnetURI{},
			Quality:    "SuperDuper5",
			Resolution: "42K",
			Sound:      "H264",
			Name:       "Keanu4Evah5",
			Seeders:    42,
			Leechers:   42,
		},
	})
	return []visual_resource.Resource{
		mockResource,
	}, nil
}