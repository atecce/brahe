//
// I should not like my writing to spare other people the trouble of thinking.
// But, if possible, to stimulate someone to thoughts of their own.
//

package main

import (
	"lyrics_net"
	"os"
)

func main() {

	// specify where to start
	start := os.Args[1]

	// start the investigation
	lyrics_net.Investigate(start)
}
