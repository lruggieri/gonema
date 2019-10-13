package rarbg

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	"gitlab.com/ruggieri/gonema/pkg/utils"
	"os"
	"path"
	"testing"
)

func TestCaptchaSelection(t *testing.T){
	mainContext := context.Background()

	ctx, cancel := chromedp.NewContext(mainContext,
		//chromedp.WithDebugf(log.Printf),
	)
	defer cancel()


	wd, err := os.Getwd()
	if err != nil{
		t.Error(err)
		t.FailNow()
	}

	//using a real RARBG HTML page as testing file
	testFilePath := "file://" + path.Join(wd,"testing_files","rarbg_search_page.html")

	fmt.Println("Navigating to test file",testFilePath)
	err = chromedp.Run(ctx,
		chromedp.Navigate(testFilePath),
	)
	if err != nil{
		t.Error(err)
		t.FailNow()
	}

	fmt.Println("Navigated")

	//if we try to detect captcha now, we would get an error
	utils.DebugActive=true
	utils.Logger.Level = logrus.DebugLevel
	err = dealWithThreatDefencePage(ctx)
	if err != context.DeadlineExceeded{
		t.Error("Expecting to get '",context.DeadlineExceeded,"' error")
		t.FailNow()
	}
}

func TestTorrentTitleSelection(t *testing.T){
	ctx, cancel := chromedp.NewContext(context.Background(),
		//chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	wd, err := os.Getwd()
	if err != nil{
		t.Error(err)
		t.FailNow()
	}

	//using a real RARBG HTML page as testing file
	testFilePath := "file://" + path.Join(wd,"testing_files","rarbg_search_page.html")
	fmt.Println("Navigating to test file",testFilePath)
	var nodes []*cdp.Node
	err = chromedp.Run(ctx,
		chromedp.Navigate(testFilePath),
	)
	if err != nil{
		t.Error(err)
		t.FailNow()
	}

	fmt.Println("Navigated")

	err = chromedp.Run(ctx,
		chromedp.Nodes(`tr[class="lista2"] > td:nth-child(2) > a:nth-child(1)`, &nodes, chromedp.ByQueryAll),
	)
	if err != nil{
		t.Error(err)
		t.FailNow()
	}

	expectedNodesLength := 25
	if len(nodes) != expectedNodesLength{
		t.Error("Expected nodes len to be ",expectedNodesLength)
		t.FailNow()
	}
}


func TestMainFilmPageInfoSelection(t *testing.T){
	scraper := Scraper{}
	ctx, cancel := chromedp.NewContext(context.Background(),
		//chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	wd, err := os.Getwd()
	if err != nil{
		t.Error(err)
		t.FailNow()
	}

	//using a real RARBG HTML page as testing file
	testFilePath := "file://" + path.Join(wd,"testing_files","rarbg_specific_torrent_page.html")
	fmt.Println("Navigating to test file",testFilePath)
	err = chromedp.Run(ctx,
		chromedp.Navigate(testFilePath),
	)
	if err != nil{
		t.Error(err)
		t.FailNow()
	}
	fmt.Println("Navigated")

	resultingTorrent, err := scraper.getSinglePageTorrentInfo(ctx,testFilePath)

	if err != nil{
		t.Error(err)
		t.FailNow()
	}

	if resultingTorrent == nil{
		t.Error("no torrent could be fetched")
		t.FailNow()
	}
	spew.Dump(resultingTorrent)
}