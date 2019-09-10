package rarbg

import (
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/davecgh/go-spew/spew"
	"os"
	"path"
	"testing"
	"context"
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
	for _, node := range nodes{
		spew.Dump(node.AttributeValue("title"))
	}
}
