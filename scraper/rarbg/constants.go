package rarbg

const ThreatDefencePageTag = "threat_defence.php"

//CAPTCHA RESOLUTION
const CaptchaPageImageBox = `/html/body/form/div/div/table/tbody/tr[2]/td[2]`
const CaptchaPageImagePath = CaptchaPageImageBox+`/img`
const CaptchaStringInputID = `#solve_string`
const CaptchaStringButtonSubmitID = `#button_submit`
