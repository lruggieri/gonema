package scraper

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"log"
	"math"
	"testing"
	"time"
)

func TestChromeDP(t *testing.T){

	ctx, cancel := chromedp.NewContext(context.Background(),/*chromedp.WithDebugf(log.Printf)*/)
	defer cancel()

	// run task list
	var buf []byte
	var bufElement []byte
	err := chromedp.Run(ctx,
		//page.SetDownloadBehavior(page.SetDownloadBehaviorBehaviorAllow).WithDownloadPath("/home/luca/go/central/src/gonema/scraper"),
		chromedp.Navigate(`https://rarbgunblock.com/threat_defence.php?defence=1&r=24320121`),
		chromedp.Sleep(6*time.Second),
		chromedp.Screenshot(`/html/body/form/div/div/table/tbody/tr[2]/td[2]/img`, &bufElement, chromedp.BySearch),
		fullScreenshot(90, &buf),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile("elementScreenshot.png", bufElement, 0644); err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}


	if err := ioutil.WriteFile("fullScreenshot.png", buf, 0644); err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}


	resourceTree, err := page.GetResourceTree().Do(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resourceTree)

}
func fullScreenshot(quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			// get layout metrics
			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			// force viewport emulation
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}

			// capture screenshot
			*res, err = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      contentSize.Y,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
}