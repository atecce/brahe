package main

import (
	"bodhi/connection"
	"bodhi/db"
	"net/url"
	"os"
	"strconv"
)

//const trackID = 5151298

// initalize api url
var api = &url.URL{
	Scheme:   "http",
	Host:     "api.soundcloud.com",
	RawQuery: "client_id=" + os.Getenv("CLIENT_ID"),
}

var canvas = &db.Canvas{
	Kind: "mysql",
	URL:  "root:@tcp(127.0.0.1:3306)/",
	Name: "canvas",
}

func main() {

	// set the canvas
	canvas.Initiate()

	// start counter at last ID
	var trackID int
	canvas.AddTable("track")
	canvas.GetLatest(&trackID, "track")
	for {

		// increment ID
		trackID++

		// attempt to get info on trackID
		api.Path = "tracks/" + strconv.Itoa(trackID)

		// try and communicate
		connection.Communicate(api, canvas)
	}
}
