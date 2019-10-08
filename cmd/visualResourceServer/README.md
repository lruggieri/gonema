#### RESOURCE INFO API

At this moment resources can be retrieved only through *imdbID*

`/resourceInfo?imdbID=tt6146586`

no local database is implemented at the moment, and the first search for each new *imdbID* can possibly
take some time to be fetched (due to captcha recognition and general crawling) but subsequent request
on the same *imdbID* should take much less time, due to caching mechanism.
Only Torrent information is returned at the moment.