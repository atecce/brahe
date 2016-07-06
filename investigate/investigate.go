package main

import (
	"investigations/lyrics_net"
	"flag"
)

func main() {

	// set start flag
	start := flag.String("s", "0", "Specify start artist of crawl.")
	flag.Parse()

	// start the investigation
	lyrics_net.Investigate(*start)
}
