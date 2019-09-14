package rarbg

//GENERAL
const Name = "RARBG"
const initialScrapersPoolSize = 10

//CAPTCHA RESOLUTION
const threatDefencePageTag = "threat_defence.php"
const captchaPageImageBox = `/html/body/form/div/div/table/tbody/tr[2]/td[2]`
const captchaPageImagePath = captchaPageImageBox +`/img`
const captchaStringInputID = `#solve_string`
const captchaStringButtonSubmitID = `#button_submit`


//MAIN TORRENT LIST PAGE
const mainTorrentListPageLinks = `tr[class="lista2"] > td:nth-child(2) > a:nth-child(1)`

//SPECIFIC TORRENT PAGE
const specificTorrentPageMagnet = `html > body > table:nth-child(6) > tbody > 
					tr:nth-child(1) > td:nth-child(2) > div > table > tbody > tr:nth-child(2) > 
					td:nth-child(1) > div > table > tbody > tr:nth-child(1) > td:nth-child(2) > a:nth-child(3)`