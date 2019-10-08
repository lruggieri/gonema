package main

import (
	"flag"
)


var visualResourceServer = flag.String("visualResourceServer", "","visualResourceServer host to use for this stress test")
var tps = flag.Int("tps",10,"amount of messages to inject to the specified visualResourceServer url PER SECOND")

func main(){
	flag.Parse()

	if len(*visualResourceServer) == 0{
		panic("please specify a valid resource server")
	}


	parametersToShoot := []string{
		"imdbID=tt6146586",
	}

}

