package rarbg

const ThreatDefencePageTag = "threat_defence.php"

//CAPTCHA RESOLUTION
const CaptchaPageImageBox = `/html/body/form/div/div/table/tbody/tr[2]/td[2]`
const CaptchaPageImagePath = CaptchaPageImageBox+`/img`
const CaptchaStringInputID = `#solve_string`
const CaptchaStringButtonSubmitID = `#button_submit`


//MAIN TORRENT LIST PAGE
const MainTorrentListPageLinks = `tr[class="lista2"] > td:nth-child(2) > a:nth-child(1)`

//SPECIFIC TORRENT PAGE
const SpecificTorrentPageMagnet = `html > body > table:nth-child(6) > tbody > 
					tr:nth-child(1) > td:nth-child(2) > div > table > tbody > tr:nth-child(2) > 
					td:nth-child(1) > div > table > tbody > tr:nth-child(1) > td:nth-child(2) > a:nth-child(3)`
