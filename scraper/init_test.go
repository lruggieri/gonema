package scraper

import (
	"fmt"
	"github.com/gocolly/colly"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestGeneral(t *testing.T){
	url := "https://rarbgunblock.com/torrents.php"

	// Instantiate default collector
	c := colly.NewCollector(
		// Attach a debugger to the collector
		/*colly.Debugger(&debug.LogDebugger{}),
		colly.Async(true),*/
	)

	c.RedirectHandler = func(req *http.Request, via []*http.Request) error{
		fmt.Println("I have been redirected!")
		return nil
	}

	c.OnHTML("a[title]", func(e *colly.HTMLElement) {
		fmt.Printf("a[title]: %s\n", e.Attr("title"))
		fmt.Printf("path: %s\n", e.Request.URL.Path)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("OnRequest: %s\n", r.URL.Path)
	})


	c.Limit(&colly.LimitRule{
		DomainGlob:  "*threat*",
		Parallelism: 3,
		RandomDelay: 5 * time.Second,
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Printf("OnResponse: %s\n", r.Request.URL.Path)

		if strings.Contains(r.Request.URL.Path,"threat_defence"){
			fmt.Println("Threat defence page found. Waiting for captcha to be rendered...")
			time.Sleep(10*time.Second)
		}
		fmt.Println("time over")
		c.OnHTML("html", func(e *colly.HTMLElement) {
			fmt.Println(e.ChildText("script"))
		})
		c.OnRequest(func(r *colly.Request) {
			fmt.Printf("OnRequest: %s\n", r.URL.Path)
		})

		c.OnResponse(func(r *colly.Response) {
			fmt.Printf("OnResponse: %s\n", r.Request.URL.Path)
		})

	})

	c.Visit(url)
	c.Wait()

}

/*
interesting URLs:


https://rarbgunblock.com/threat_defence.php?defence=1&r=24320121

*/