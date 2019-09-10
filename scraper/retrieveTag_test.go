package scraper

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRetrieveTagOLD(t *testing.T){
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ts := httptest.NewServer(writeHTML(`
<body>
<p id="content" onclick="changeText()">Original content.</p>
<script>
function changeText() {
	document.getElementById("content").textContent = "New content!"
}
</script>
</body>
	`))
	defer ts.Close()

	var outerBefore, outerAfter string
	if err := chromedp.Run(ctx,
		chromedp.Navigate(ts.URL),
		chromedp.OuterHTML("#content", &outerBefore),
		chromedp.Click("#content", chromedp.ByID),
		chromedp.OuterHTML("#content", &outerAfter),
	); err != nil {
		panic(err)
	}
	fmt.Println("OuterHTML before clicking:")
	fmt.Println(outerBefore)
	fmt.Println("OuterHTML after clicking:")
	fmt.Println(outerAfter)


}


func TestRetrieveTag(t *testing.T){
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ts := httptest.NewServer(writeHTML(
		`
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
<meta name="theme-color" content="#3860BB" />
</head>
<body>
				
<style type="text/css">a,abbr,acronym,address,applet,article,aside,audio,b,big,blockquote,body,canvas,caption,center,cite,code,dd,del,details,dfn,div,dl,dt,em,fieldset,figcaption,figure,footer,form,h1,h2,h3,h4,h5,h6,header,hgroup,html,i,iframe,img,ins,kbd,label,legend,li,mark,menu,nav,object,ol,p,pre,q,s,samp,section,small,span,strike,strong,sub,summary,sup,table,tbody,td,tfoot,th,thead,time,tr,tt,u,ul,var,video{margin:0;padding:0;border:0;outline:0;font:inherit;vertical-align:baseline}article,aside,details,figcaption,figure,footer,header,hgroup,menu,nav,section{display:block}body{line-height:1}ol,ul{list-style:none}blockquote,q{quotes:none}blockquote:after,blockquote:before,q:after,q:before{content:'';content:none}ins{text-decoration:none}del{text-decoration:line-through}table{border-collapse:collapse;border-spacing:0}
body {
    background: #000 url("https://dyncdn.me/static/20/img/bknd_body.jpg") repeat-x scroll 0 0 !important;
    font: 400 8pt normal Tahoma,Verdana,Arial,Arial  !important;
}
.button {
    background-color: #3860bb;
    border: none;
    color: white;
    padding: 15px 32px;
    text-align: center;
    text-decoration: none;
    display: inline-block;
    font-size: 16px;
    cursor: pointer;
    text-transform: none;
    overflow: visible;
}
.content-rounded {
    background: #fff none repeat scroll 0 0 !important;
    border-radius: 3px;
    color: #000 !important;
    padding: 20px;
    width:961px;
}
</style>
<script type="text/javascript" src="https://dyncdn.me/static/20/js/jquery-1.11.3.min.js"></script>
<form action="/threat_defence.php" method="GET" autocomplete="off">
<input type="hidden" name="defence" value="2">
<input type="hidden" name="sk" value="sztja9nw7b">
<input type="hidden" name="cid" value="22948950">
<input type="hidden" name="i" value="1333398476">
<input type="hidden" name="ref_cookie" value="rarbgunblock.com">
<input type="hidden" name="r" value="85665079">
<div align="center" style="margin-top:20px;padding-top:20px;color: #000 !important;">
<div  class="content-rounded" style="color: #000 !important;">
<table width="100%">
<tr><td colspan="2">
<table>
<tr>
<td valign="top" style="vertical-align:top !important;"><img src="https://dyncdn.me/static/20/img/logo_dark_nodomain2_optimized.png"></td>
<td valign="middle" style="vertical-align:middle !important;"><font color="#FF0000"><b>Our system has detected abnormal activity from your ip address 88.189.115.220</b></font><br/>

</td>
</tr>
</table>
<hr/></td></tr>
<tr>
<td valign="top" width="50%" style="vertical-align:top !important;">
<b style="font-weight:bold;">Possible reason why this happend :</b><br/>
<ul type="circle" style="list-style-type: circle !important;">
<li>You dont have javascript or cookies enabled</li>
<li>You are using a broken proxy/mirror</li>
<li>You have extensions in your browser that are blocking cookies/javascript</li>
<li>You are using a VPN that is abused on our site</li>
<li>Your pc is acting as a proxy to other networks</li>
<li>Your ip is listed in ProjectHoneypot</li>
<li>You are making more than 10 requests per second to the site</li>


</ul>
</td>
<td valign="top" width="50%" style="vertical-align:top !important;">

We are sorry but we were unable to automatically verify your browser!<br/> <b style="font-weight:bold !important;">Please enter the captcha below to continue</b><br/>
<img src="/threat_captcha.php?cid=22948950_yf3z6_1333398476&r=56154826" lazyload="off" />
<br/>
<input type="text" name="solve_string" id="solve_string" value="" autocomplete="off" maxlength="5" placeholder="Enter Captcha" style="width:160px !important;" />
<input type="hidden" name="captcha_id" value="z92euwsdoy8k635bqnxrhm4fcgi7atvp" />
<input type="hidden" name="submitted_bot_captcha" value="1" />
<br/>

<button class="button" type="submit" style="display:none;padding-top:10px;" id="button_submit"><i class="icon-user"></i> I am human</button>

</td>
</tr>
</table>






</form>
<hr/>
<table width="959" style="width:959px !important;">
<tr>
<td width="50%">
<h2 style="font-size: 20px;line-height: 1.3;font-weight: 400;color: #404040;margin: 0;padding: 0;">Why do I have to complete a CAPTCHA?</h2>
<p style="margin-top: 1em;color: #404040;">Completing the CAPTCHA proves you are a human and gives you temporary access to this website.</p>
</td>
<td width="50%">
<h2 style="font-size: 20px;line-height: 1.3;font-weight: 400;color: #404040;margin: 0;padding: 0;">What can I do to prevent this in the future?</h2>
<p style="margin-top: 1em;color: #404040;">If you are on a personal connection, like at home, you can run an anti-virus scan on your device to make sure it is not infected with malware.<br/><br/>
If you are at an office or shared network, you can ask the network administrator to run a scan across the network looking for misconfigured or infected devices.</p>
</td>
</tr>
</table>
</div>

</div>


</body>
	`))
	defer ts.Close()

	var attrs []map[string]string
	if err := chromedp.Run(ctx,
		chromedp.Navigate(ts.URL),
		chromedp.AttributesAll("img", &attrs, chromedp.ByQueryAll),
	); err != nil {
		panic(err)
	}
	fmt.Println("OuterHTML before clicking:")
	fmt.Println(attrs)


}

func writeHTML(content string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		io.WriteString(w, strings.TrimSpace(content))
	})
}