package main

import (
	"flag"
	"investigations/websites"
)

func main() {

	// set start flag
	start := flag.String("s", "0", "Specify start artist of crawl.")
	flag.Parse()

	// start the investigation
	websites.Investigate(*start)
}
