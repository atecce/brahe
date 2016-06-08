package main

import (
	"flag"
	"lyrics_net"
)

func main() {

	// set start flag
	start := flag.String("s", "0", "Specify start artist of crawl.")
	flag.Parse()

	// start the investigation
	lyrics_net.Investigate(*start)
}
