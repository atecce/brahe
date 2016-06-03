//
// I should not like my writing to spare other people the trouble of thinking.
// But, if possible, to stimulate someone to thoughts of their own.
//

package main

import (
	"flag"
	"lyrics_net"
)

func main() {

	// set start flag
	start := flag.String("s", "0", "Specify start of crawl.")
	flag.Parse()

	// start the investigation
	lyrics_net.Investigate(*start)
}
