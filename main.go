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

func main() {

	// set the canvas
	canvas := db.Initiate()
	defer canvas.Close()

	// start counter at last ID
	var trackID int
	db.AddTable("track", canvas)
	db.GetLatest(&trackID, "track", canvas)
	for {

		// increment ID
		trackID++

		// attempt to get info on trackID
		api.Path = "tracks/" + strconv.Itoa(trackID)

		// try and communicate
		connection.Communicate(api, canvas)
	}
}
