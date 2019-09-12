package rarbg

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"os"
	"path"
	"testing"
)

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
	var nodes []*cdp.Node
	err = chromedp.Run(ctx,
		chromedp.Nodes(
			`html > body > table:nth-child(3) > tbody > 
					tr:nth-child(1) > td:nth-child(2) > div > table > tbody > tr:nth-child(2) > 
					td:nth-child(1) > div > table > tbody > tr:nth-child(1) > td:nth-child(2) > a:nth-child(3)`, &nodes, chromedp.ByQueryAll,
		),
	)

	expectedNodesLength := 1
	if len(nodes) != expectedNodesLength{
		t.Error("Expected nodes len to be ",expectedNodesLength)
		t.FailNow()
	}
	fmt.Println(nodes[0].AttributeValue("href"))
}