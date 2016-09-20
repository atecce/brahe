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
		Name: "canvas",
	},
}

var resources = map[string][]string{
	// "users":     {"tracks", "playlists", "followings", "followers", "comments", "favorites"},
	// "tracks":    {"comments", "favoriters"},
	"playlists": nil,
	"comments":  nil,
}

var wg sync.WaitGroup

func main() {

	// set the canvas TODO maybe close the clients
	// api.Canvas.Initiate()
	// defer api.Canvas.Session.Close()
	for id := 0; ; id++ {
		for resource := range resources {

			subresources := resources[resource]

			// api.Canvas.AddFamily(resource)

			// check for ids already present
			// missing := api.Canvas.GetMissing(table)
			// log.Println(missing)

			methodPath := resource + "/" + strconv.Itoa(id)

			// construct method
			method := &url.URL{
				Scheme:   "http",
				Host:     "api.soundcloud.com",
				Path:     methodPath,
				RawQuery: "client_id=" + os.Getenv("CLIENT_ID"),
			}

			// try and communicate
			if skip := api.Communicate(resource, method); skip {
				continue
			}

			for _, subresource := range subresources {

				// populate tables concurrently
				// wg.Add(1)
				// go func(resource, subresource string) {
				// 	defer wg.Done()

				// input entries we know about
				// if _, ok := missing[id]; !ok {

				// check each subresource
				method.Path = methodPath + "/" + subresource
				api.Communicate(resource, method)
				// }
			}
			// }(resource, subresource)
		}
		println()
	}

	// wg.Wait()
}
