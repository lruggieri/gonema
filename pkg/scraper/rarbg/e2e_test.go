package rarbg

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/nmmh/magneturi/magneturi"
	"github.com/sirupsen/logrus"
	"gitlab.com/ruggieri/gonema/pkg/torrent"
	"gitlab.com/ruggieri/gonema/pkg/utils"
	"testing"
)

func TestEndToEnd(t *testing.T) {
	utils.DebugActive = true
	utils.Logger.Level = logrus.DebugLevel

	scraper := Scraper{}

	type test struct{
		imdbID string
		expectedTorrents []torrent.Torrent
	}

	tests := []test{
		{
			"tt6146584",
			nil,
		},
		{
			"tt6146586",
			[]torrent.Torrent{
				{MagnetLink: getMagnetUri("magnet:?xt=urn:btih:6ed79bc07e02064b35fdc8c5858810e9601c5ef2&dn=John.Wick.Chapter.3.Parabellum.2019.2160p.UHD.BluRay.X265.10bit.HDR.TrueHD.7.1.Atmos-TERMiNAL&tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&tr=udp%3A%2F%2F9.rarbg.me%3A2730&tr=udp%3A%2F%2F9.rarbg.to%3A2760"),},
				{MagnetLink: getMagnetUri("magnet:?xt=urn:btih:22490400fd53616e4f53999030675b68457021c0&dn=John.Wick.Chapter.3.Parabellum.2019.BDRip.x264-SPARKS&tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&tr=udp%3A%2F%2F9.rarbg.me%3A2720&tr=udp%3A%2F%2F9.rarbg.to%3A2770"),},
				{MagnetLink: getMagnetUri("magnet:?xt=urn:btih:46812fc999e1052d0a1aa2e3fcd9f58abdadf0fb&dn=John.Wick.Chapter.3.Parabellum.2019.1080p.BluRay.REMUX.AVC.DTS-HD.MA.TrueHD.7.1.Atmos-FGT&tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&tr=udp%3A%2F%2F9.rarbg.me%3A2740&tr=udp%3A%2F%2F9.rarbg.to%3A2720"),},
				{MagnetLink: getMagnetUri("magnet:?xt=urn:btih:dd4e5ae4e0ff7c144249198d202f13dee1bfed4f&dn=John.Wick.Chapter.3.Parabellum.2019.1080p.BluRay.AVC.TrueHD.7.1.Atmos-HDC&tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&tr=udp%3A%2F%2F9.rarbg.me%3A2760&tr=udp%3A%2F%2F9.rarbg.to%3A2770"),},
				{MagnetLink: getMagnetUri("magnet:?xt=urn:btih:3bd6a308d85ea56167577d65bef800a6b8d0399a&dn=John.Wick.Chapter.3.Parabellum.2019.1080p.BluRay.x264-SPARKS&tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&tr=udp%3A%2F%2F9.rarbg.me%3A2730&tr=udp%3A%2F%2F9.rarbg.to%3A2740"),},
				{MagnetLink: getMagnetUri("magnet:?xt=urn:btih:3ba96b506d52cd0378a240ab91d5bd0f923c9bd1&dn=John.Wick.Chapter.3.Parabellum.2019.720p.BluRay.x264-SPARKS&tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&tr=udp%3A%2F%2F9.rarbg.me%3A2760&tr=udp%3A%2F%2F9.rarbg.to%3A2740"),},
				{MagnetLink: getMagnetUri("magnet:?xt=urn:btih:40e252e723e636b9896c02e2c38a362da8f7699b&dn=John.Wick.Chapter.3.Parabellum.2019.2160p.BluRay.x265.10bit.SDR.DTS-HD.MA.TrueHD.7.1.Atmos-SWTYBLZ&tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&tr=udp%3A%2F%2F9.rarbg.me%3A2720&tr=udp%3A%2F%2F9.rarbg.to%3A2780"),},
				{MagnetLink: getMagnetUri("magnet:?xt=urn:btih:bc000e05f7beabc8921805ad169e10ce9302a0e4&dn=John.Wick.Chapter.3.Parabellum.2019.2160p.BluRay.x265.10bit.HDR.DTS-HD.MA.TrueHD.7.1.Atmos-SWTYBLZ&tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&tr=udp%3A%2F%2F9.rarbg.me%3A2770&tr=udp%3A%2F%2F9.rarbg.to%3A2780"),},
				{MagnetLink: getMagnetUri("magnet:?xt=urn:btih:0c03d49d23cdb94da78bf58c1b4a0397b581938e&dn=John.Wick.Chapter.3.Parabellum.2019.2160p.BluRay.x264.8bit.SDR.DTS-HD.MA.TrueHD.7.1.Atmos-SWTYBLZ&tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&tr=udp%3A%2F%2F9.rarbg.me%3A2740&tr=udp%3A%2F%2F9.rarbg.to%3A2760"),},
				{MagnetLink: getMagnetUri("magnet:?xt=urn:btih:38d65d6cbd61c55c762a3f2acf8ca9d851b05fe3&dn=John.Wick.Chapter.3.Parabellum.2019.2160p.BluRay.REMUX.HEVC.DTS-HD.MA.TrueHD.7.1.Atmos-FGT&tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&tr=udp%3A%2F%2F9.rarbg.me%3A2780&tr=udp%3A%2F%2F9.rarbg.to%3A2710"),},
				{MagnetLink: getMagnetUri("magnet:?xt=urn:btih:8bb08c43f4c3b161f2b1ca51ddf499eb918caee8&dn=John.Wick.Chapter.3.Parabellum.2019.2160p.BluRay.HEVC.TrueHD.7.1.Atmos-BHD&tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&tr=udp%3A%2F%2F9.rarbg.me%3A2720&tr=udp%3A%2F%2F9.rarbg.to%3A2710"),},
				{MagnetLink: getMagnetUri("magnet:?xt=urn:btih:d1c58426f6af9b7310bd82c2930b9f287b207536&dn=John.Wick.Chapter.3.Parabellum.2019.1080p.BluRay.H264.AAC-RARBG&tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&tr=udp%3A%2F%2F9.rarbg.me%3A2730&tr=udp%3A%2F%2F9.rarbg.to%3A2750"),},
				{MagnetLink: getMagnetUri("magnet:?xt=urn:btih:41b15c9c7ce0db59043d26221e4180a5664e39ec&dn=John.Wick.Chapter.3.Parabellum.2019.1080p.BluRay.x264.TrueHD.7.1.Atmos-HDC&tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&tr=udp%3A%2F%2F9.rarbg.me%3A2740&tr=udp%3A%2F%2F9.rarbg.to%3A2770"),},
				{MagnetLink: getMagnetUri("magnet:?xt=urn:btih:63738f120ecc3f5f5f832742ce88e1518a2fa6d5&dn=John.Wick.Chapter.3.Parabellum.2019.1080p.BluRay.x264.DTS-HDC&tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&tr=udp%3A%2F%2F9.rarbg.me%3A2770&tr=udp%3A%2F%2F9.rarbg.to%3A2760"),},
				{MagnetLink: getMagnetUri("magnet:?xt=urn:btih:79ba705c58ac46e83881c65fbcd90befc2300e8e&dn=John.Wick.Chapter.3.Parabellum.2019.1080p.BluRay.x264.DTS-HD.MA.7.1-HDC&tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&tr=udp%3A%2F%2F9.rarbg.me%3A2790&tr=udp%3A%2F%2F9.rarbg.to%3A2760"),},
				{MagnetLink: getMagnetUri("magnet:?xt=urn:btih:62d5a098f67b151ba1a023d9c739c61643119c9b&dn=John.Wick.Chapter.3.Parabellum.2019.720p.BluRay.H264.AAC-RARBG&tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&tr=udp%3A%2F%2F9.rarbg.me%3A2730&tr=udp%3A%2F%2F9.rarbg.to%3A2760"),},
				{MagnetLink: getMagnetUri("magnet:?xt=urn:btih:ae955c3ad0395e1d2cd70bc29f8a8957b6131957&dn=John.Wick.Chapter.3.Parabellum.2019.BRRip.XviD.AC3-XVID&tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&tr=udp%3A%2F%2F9.rarbg.me%3A2770&tr=udp%3A%2F%2F9.rarbg.to%3A2750"),},
				{MagnetLink: getMagnetUri("magnet:?xt=urn:btih:e1bbd3679127694d9495824423e02083ed706bc2&dn=John.Wick.Chapter.3.Parabellum.2019.720p.BRRip.XviD.AC3-XVID&tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&tr=udp%3A%2F%2F9.rarbg.me%3A2760&tr=udp%3A%2F%2F9.rarbg.to%3A2750"),},
				{MagnetLink: getMagnetUri("magnet:?xt=urn:btih:fb49248abd00d73c8732eb3a09e3b85d3562f648&dn=John.Wick.Chapter.3.Parabellum.2019.BRRip.XviD.MP3-XVID&tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&tr=udp%3A%2F%2F9.rarbg.me%3A2790&tr=udp%3A%2F%2F9.rarbg.to%3A2740"),},
				{MagnetLink: getMagnetUri("magnet:?xt=urn:btih:0d91c1387a094c84c55bf4ae276e3213923e0adb&dn=John.Wick.Chapter.3.Parabellum.2019.720p.BluRay.x264.DD5.1-HDC&tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&tr=udp%3A%2F%2F9.rarbg.me%3A2730&tr=udp%3A%2F%2F9.rarbg.to%3A2740"),},
				{MagnetLink: getMagnetUri("magnet:?xt=urn:btih:8af5bedd19001bd37329f9587c9f28afce0e4831&dn=John.Wick.Chapter.3.Parabellum.2019.1080p.WEBRip.x264-RARBG&tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&tr=udp%3A%2F%2F9.rarbg.me%3A2760&tr=udp%3A%2F%2F9.rarbg.to%3A2790"),},
				{MagnetLink: getMagnetUri("magnet:?xt=urn:btih:69fa152febefbef225ba1ac2ded95dca01be6f38&dn=John.Wick.Chapter.3.Parabellum.2019.WEBRip.x264-ION10&tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&tr=udp%3A%2F%2F9.rarbg.me%3A2750&tr=udp%3A%2F%2F9.rarbg.to%3A2740"),},
				{MagnetLink: getMagnetUri("magnet:?xt=urn:btih:d4e0bd1c7071c2219e32bccbe5ebf1c834350d6f&dn=John.Wick.Chapter.3.Parabellum.2019.WEBRip.XviD.MP3-FGT&tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&tr=udp%3A%2F%2F9.rarbg.me%3A2780&tr=udp%3A%2F%2F9.rarbg.to%3A2730"),},
				{MagnetLink: getMagnetUri("magnet:?xt=urn:btih:2ce226de41c8cd8610fdbb0450c50fc24b8d96c8&dn=John.Wick.Chapter.3.Parabellum.2019.WEBRip.XviD.AC3-FGT&tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&tr=udp%3A%2F%2F9.rarbg.me%3A2740&tr=udp%3A%2F%2F9.rarbg.to%3A2790"),},
				{MagnetLink: getMagnetUri("magnet:?xt=urn:btih:a4fd174f8e8454a0dfe3162ab72bf680bdbca626&dn=John.Wick.Chapter.3.Parabellum.2019.720p.WEBRip.XviD.AC3-FGT&tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&tr=udp%3A%2F%2F9.rarbg.me%3A2730&tr=udp%3A%2F%2F9.rarbg.to%3A2750"),},
			},
		},
	}

	for tstIdx,tst := range tests{
		finalTorrents,err := scraper.GetTorrentLinks("",tst.imdbID)
		if err != nil{
			t.Error(err)
			t.FailNow()
		}
		if tst.expectedTorrents == nil && finalTorrents != nil{
			t.Error("Test #",tstIdx," was expecting nil as result, but got something else")
			t.FailNow()
		}
		if len(tst.expectedTorrents) != len(finalTorrents){
			spew.Dump(finalTorrents)
			t.Error("Test #",tstIdx," was expecting results with length ",len(tst.expectedTorrents),
				" but got ",len(finalTorrents))
			t.FailNow()
		}
		for expectedTorrentIdx,expectedTorrent := range tst.expectedTorrents{
			resTorrent := finalTorrents[expectedTorrentIdx]

			resMagnetLink := applyMainParametersFilter(resTorrent.MagnetLink)
			if expectedTorrent.MagnetLink.String() != resMagnetLink.String(){
				t.Error("Test #",tstIdx," was expecting result torrent #",expectedTorrentIdx," to be " +
					"'",expectedTorrent.MagnetLink.String(),"' but got '",resMagnetLink.String(),"'")
				t.FailNow()
			}
		}

	}
}


func getMagnetUri(iMagnetLink string)magneturi.MagnetURI{
	ret,_ := magneturi.Parse(iMagnetLink,false)
	return applyMainParametersFilter(*ret)
}
//filtering the value that changes continually, like 'tr'
func applyMainParametersFilter(iMagnet magneturi.MagnetURI) magneturi.MagnetURI{
	magnetMainParameters,_ := iMagnet.Filter("xt","dn")
	return *magnetMainParameters
}