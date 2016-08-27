package main

import (
	"bodhi/connection"
	"bodhi/db"
	"net/url"
	"os"
	"strconv"
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

func main() {

	// set the canvas
	api.Canvas.Initiate()

	// start counter at last ID
	var trackID int
	api.Canvas.AddTable("track")
	api.Canvas.GetLatest(&trackID, "track")
	for {

		// increment ID
		trackID++

		// attempt to get info on trackID
		api.Method.Path = "tracks/" + strconv.Itoa(trackID)

		// try and communicate
		api.Communicate()
	}
}
