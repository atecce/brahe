package main

import (
	"net/url"
	"os"
	"strconv"
	"sync"

	"github.com/atecce/brahe/connection"
	"github.com/atecce/brahe/db"
)

//const trackID = 5151298

var api = &connection.API{

	Canvas: &db.Canvas{
		Kind: "mysql",
		IP:   "127.0.0.1",
		Name: "canvas",
	},
}

var families = []string{
	"user",
	"track",
	"playlist",
	"comment",
}

var wg sync.WaitGroup

func main() {

	// set the canvas TODO maybe close the clients
	api.Canvas.Initiate()
	// defer api.Canvas.Session.Close()
	for _, family := range families {

		api.Canvas.AddFamily(family)

		// check for ids already present
		// missing := api.Canvas.GetMissing(table)
		// log.Println(missing)

		// populate tables concurrently
		wg.Add(1)
		go func(family string) {
			defer wg.Done()

			// input entries we know about
			for id := 0; ; id++ {
				// if _, ok := missing[id]; !ok {

				// attempt to get info on trackID
				method := &url.URL{
					Scheme:   "http",
					Host:     "api.soundcloud.com",
					Path:     family + "s/" + strconv.Itoa(id),
					RawQuery: "client_id=" + os.Getenv("CLIENT_ID"),
				}

				// try and communicate
				api.Communicate(family, method)
				// }
			}
		}(family)
	}

	wg.Wait()
}
