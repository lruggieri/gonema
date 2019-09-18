package rarbg

import (
	"context"
	"errors"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/nmmh/magneturi/magneturi"
	"github.com/otiai10/gosseract"
	"gitlab.com/ruggieri/gonema/pkg/torrent"
	"gitlab.com/ruggieri/gonema/pkg/utils"
	"io/ioutil"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

var availableScrapers []context.Context
var availableScrapersLock sync.RWMutex

//initialize a pool of RARBG scrapers
func init(){
	availableScrapers = make([]context.Context, 0, initialScrapersPoolSize)
	for i := 0 ; i < cap(availableScrapers) ; i++{
		newBaseContext := context.Background()
		//we don't want to ever cancel a context, so let's not return it's cancel function
		ctx, _ := chromedp.NewContext(newBaseContext,
			//chromedp.WithDebugf(log.Printf),
		)

		availableScrapers = append(availableScrapers, ctx)
	}
}
//returns a random scraper, eventually re-initializing it if non valid
func getScraperRoundRobin() context.Context{
	availableScrapersLock.Lock()
	defer availableScrapersLock.Unlock()

	randomScraperPosition := utils.GetRandomPositiveInt(len(availableScrapers))


	if utils.DebugActive{utils.Logger.Debug("Using random scraper #"+strconv.Itoa(randomScraperPosition))}
	randomScraper := availableScrapers[randomScraperPosition]
	if randomScraper.Err() != nil{
		if utils.DebugActive{utils.Logger.Debug("Re-initializing context")}
		availableScrapers[randomScraperPosition] = context.Background()
	}

	return randomScraper
}


//Main RARBG scraper type
type Scraper struct{}

/*
TODO implement name usage ('/torrents.php?search=...')
*/
func (sc *Scraper)GetTorrentLinks(iResourceName, iResourceImdbID string) (oTorrents []torrent.Torrent, oErr error){
	mainDomain := "https://rarbgunblock.com"
	mainSearchLink := mainDomain+"/torrents.php?imdb="+ iResourceImdbID

	if utils.DebugActive{utils.Logger.Debug("Creating new context")}
	//we don't want to ever cancel a context, so let's not return it's cancel function
	randMainCtx := getScraperRoundRobin()


	err := navigateWithCaptchaDetection(randMainCtx,mainSearchLink)
	if err != nil{
		return nil,err
	}

	//getting the full list film nodes info
	specificTorrentNodesToCrawl := make([]*cdp.Node,0)
	timeout, err := executeRunWithTimeout(
		randMainCtx,
		time.Second,
		chromedp.Nodes(mainTorrentListPageLinks, &specificTorrentNodesToCrawl, chromedp.ByQueryAll),
	)
	if timeout{
		//there is nothing to be fetched in this page
		return nil,nil
	}
	if err != nil{
		return nil,err
	}

	specificTorrentLinksToCrawl := make([]string,len(specificTorrentNodesToCrawl))
	for torrentNodeIdx, torrentNode := range specificTorrentNodesToCrawl{
		specificTorrentLinksToCrawl[torrentNodeIdx] = torrentNode.AttributeValue("href")
	}

	finalTorrents := make([]torrent.Torrent,0,len(specificTorrentLinksToCrawl))
	//TODO possibly use goroutines. Unsure about the chance of getting banned for too much speed though.
	for _,specificTorrentLinkToCrawl := range specificTorrentLinksToCrawl{
		specificTorrentPage := mainDomain+specificTorrentLinkToCrawl
		err = navigateWithCaptchaDetection(randMainCtx, specificTorrentPage)
		if err != nil{
			return nil,err
		}
		magnetNodes := make([]*cdp.Node,0)
		timeout,err = executeRunWithTimeout(randMainCtx,
			1*time.Second,
			/*chromedp.ActionFunc(func(ctx context.Context) error {

				//
			}),*/
			chromedp.Nodes(specificTorrentPageMagnet, &magnetNodes, chromedp.ByQueryAll),
		)
		if timeout{
			//TODO log better with Logstash
			utils.Logger.Error(errors.New("timeout when fetching data for torrent page "+specificTorrentPage+ ", continuing..."))
			continue
		}
		if err != nil{
			return nil,err
		}
		magnetLinkString := magnetNodes[0].AttributeValue("href")
		magnetLink, err := magneturi.Parse(magnetLinkString,false)
		if err != nil{
			return nil,err
		}

		finalTorrents = append(finalTorrents, torrent.Torrent{
			MagnetLink:*magnetLink, //let's bring along pointers when not needed please! Have mercy for the heap!
		})
	}


	return finalTorrents,nil
}


/*
Remember to pass cookies in the format key1,value1,...keyN,valueN. So they must be an even number
*/
func setCookies(cookies ...string) chromedp.ActionFunc{
	return chromedp.ActionFunc(func(ctx context.Context) error {
		// create cookie expiration
		expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))
		// add cookies to chrome
		for i := 0; i < len(cookies); i += 2 {
			success, err := network.SetCookie(cookies[i], cookies[i+1]).
				WithDomain(".rarbgunblock.com").
				WithExpires(&expr).
				Do(ctx)
			if err != nil {
				return err
			}
			if !success {
				return fmt.Errorf("could not set cookie %q to %q", cookies[i], cookies[i+1])
			}
		}
		return nil
	})
}
func navigateWithCaptchaDetection(iCtx context.Context, iTargetPage string) error{
	//try to get to the main page, possibly dealing with threat security pages, for a maximum amount of time
	const maxMainPageTentatives = 3
	currentMainPageTentatives := 0

	/*
	Even after decoding the captcha, if necessary, we get redirected with the main page, and not to the page search initially (with
	the iMDB film ID), so we need a round of navigation even after decoding the captcha
	*/
	for{

		if utils.DebugActive{utils.Logger.Debug("Navigating to "+ iTargetPage +", tentative "+strconv.Itoa(currentMainPageTentatives))}
		err := chromedp.Run(iCtx,
			//page.SetDownloadBehavior(page.SetDownloadBehaviorBehaviorAllow).WithDownloadPath("/home/luca/go/central/src/gonema/scraper"),
			chromedp.Navigate(iTargetPage),
		)
		if err != nil{
			return err
		}

		//print cookies
		/*if utils.DebugActive{
			err = chromedp.Run(iCtx,
				chromedp.ActionFunc(func(ctx context.Context) error {
					cookies, err := network.GetAllCookies().Do(ctx)
					if err != nil {
						return err
					}

					for i, cookie := range cookies {
						log.Printf("chrome cookie %d: %+v", i, cookie)
					}

					return nil
				}),
			)
			if err != nil{
				return err
			}
		}*/

		landedOnExpectedPage, landedOnThreatDefencePage, err := landedRARBGPageInfo(iTargetPage,iCtx)
		if err != nil{
			return err
		}
		if landedOnExpectedPage{
			if utils.DebugActive{utils.Logger.Debug("landed on expected page: "+ iTargetPage)}
			break
		}
		currentMainPageTentatives ++
		if currentMainPageTentatives > maxMainPageTentatives{
			return errors.New("cannot land on page "+ iTargetPage +". Max tentatives ("+strconv.Itoa(maxMainPageTentatives)+") reached")
		}

		if landedOnThreatDefencePage{
			if utils.DebugActive{utils.Logger.Debug("Threat defence page triggered")}
			err = dealWithThreatDefencePage(iCtx)
			if err != nil{
				return err
			}
		}
	}
	return nil
}
func executeRunWithTimeout(iFatherCtx context.Context, iTimeoutDuration time.Duration, iActions ...chromedp.Action) (oTimeout bool, oErr error){
	newChildContext,_ := context.WithDeadline(iFatherCtx, time.Now().Add(iTimeoutDuration))

	err := chromedp.Run(newChildContext,iActions...)
	if err != nil{
		if err == context.DeadlineExceeded{
			return true, err
		}
		return false, err
	}

	return false, nil
}
func dealWithThreatDefencePage(iParentCtx context.Context) (oErr error){
	maxCaptchaCheckTrials := 3
	var threatCaptchaBeforeResolutionImageBytes []byte //page screenshot before resolving captcha
	var threatCaptchaBoxAfterResolutionBytes []byte //page screenshot after resolving captcha
	var threatCaptchaBoxAfterRedirectionBytes []byte //page screenshot after clicking on submit captcha and being redirected



	captchaFound := false
	captchaPageWaitTime := 10 * time.Second

	newChildContext,cancel := context.WithDeadline(iParentCtx,time.Now().Add(captchaPageWaitTime))
	defer cancel()

	for i := 0 ; i < maxCaptchaCheckTrials ; i++{
		if utils.DebugActive{utils.Logger.Debug("Trying to take a screenShot of the captcha to be decoded ... " +
			"tentative "+strconv.Itoa(i+1)+"/"+strconv.Itoa(maxCaptchaCheckTrials))}

		err := chromedp.Run(newChildContext,
			chromedp.WaitVisible(captchaPageImagePath,chromedp.BySearch),
			chromedp.Sleep(500*time.Millisecond), //necessary for some reason. Without, the following screenshot fails
			chromedp.Screenshot(
				captchaPageImagePath,
				&threatCaptchaBeforeResolutionImageBytes,
				chromedp.BySearch),
		)
		if err != nil{
			return err
		}
		if threatCaptchaBeforeResolutionImageBytes == nil || len(threatCaptchaBeforeResolutionImageBytes) == 0{
			continue
		}

		if utils.DebugActive{utils.Logger.Debug("Captcha screen taken")}
		captchaFound = true


		if utils.DebugActive{utils.Logger.Debug("Calling tesseract to decode image")}
		client := gosseract.NewClient()
		defer client.Close()
		client.SetImageFromBytes(threatCaptchaBeforeResolutionImageBytes)
		decodedCaptcha, _ := client.Text()
		if utils.DebugActive{utils.Logger.Debug("image decoded. result: "+decodedCaptcha)}

		err = chromedp.Run(newChildContext,
			chromedp.SendKeys(captchaStringInputID,decodedCaptcha,chromedp.ByID),
			chromedp.Screenshot(
				captchaPageImageBox,
				&threatCaptchaBoxAfterResolutionBytes,
				chromedp.BySearch),
			chromedp.Click(captchaStringButtonSubmitID,chromedp.ByID),
			chromedp.WaitVisible(mainTorrentListPageSearchBarID,chromedp.ByID), //waiting to get redirected to main page
			fullScreenShot(90, &threatCaptchaBoxAfterRedirectionBytes),         //here we should get to the main torrent list page
		)

		if utils.DebugActive{
			if err := ioutil.WriteFile("captcha.png", threatCaptchaBeforeResolutionImageBytes, 0644); err != nil {
				utils.Logger.Error(err)
			}
			if err := ioutil.WriteFile("box1.png", threatCaptchaBoxAfterResolutionBytes, 0644); err != nil {
				log.Fatal(err)
			}
			if err := ioutil.WriteFile("box2.png", threatCaptchaBoxAfterRedirectionBytes, 0644); err != nil {
				log.Fatal(err)
			}
		}

		break
	}

	if !captchaFound{
		return errors.New("captcha could not be decoded")
	}

	if utils.DebugActive{utils.Logger.Debug("threat defence page has been dealt with")}

	return nil
}
func landedRARBGPageInfo(iDesiredLink string, iCtx context.Context) (oLandedCorrectly, oIsThreatDefensePage bool, oErr error){
	targets,err := chromedp.Targets(iCtx)
	if err != nil{
		return false, false, err
	}
	for _, t := range targets{
		if strings.Contains(t.URL, threatDefencePageTag){
			return false, true, nil
		}
		if strings.Contains(t.URL, iDesiredLink){
			return true, false, nil
		}
	}

	return false, false, errors.New("landed in unknown page")
}
func fullScreenShot(quality int64, res *[]byte) chromedp.Tasks {
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