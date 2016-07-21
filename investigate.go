package main

import (
	"investigations/lyricsdotnet"
	"flag"
)

func main() {

	// set start flag
	start := flag.String("s", "0", "Specify start artist of crawl.")
	flag.Parse()

	// start the investigation
	lyricsdotnet.Investigate(*start)
}
