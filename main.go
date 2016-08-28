package main

import (
	"bodhi/connection"
	"bodhi/db"
	"net/url"
	"os"
	"strconv"
	"sync"
)

//const trackID = 5151298

var api = &connection.API{

	Canvas: &db.Canvas{
		Kind: "mysql",
		URL:  "root:@tcp(127.0.0.1:3306)/",
		Name: "canvas",
	},

	Method: &url.URL{
		Scheme:   "http",
		Host:     "api.soundcloud.com",
		RawQuery: "client_id=" + os.Getenv("CLIENT_ID"),
	},
}

var tables = []string{
	"user",
	"track",
	"playlist",
	"comment",
}

var wg sync.WaitGroup

func main() {

	// set the canvas
	api.Canvas.Initiate()

	for _, table := range tables {

		wg.Add(1)
		go func(table string) {
			defer wg.Done()

			// start counter at last ID
			for id := api.Canvas.GetLatest(table); ; id++ {

				// attempt to get info on trackID
				api.Method.Path = table + "s/" + strconv.Itoa(id)

				// try and communicate
				api.Communicate()
			}
		}(table)
	}

	wg.Wait()
}
