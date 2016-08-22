package main

import (
	"database/sql"
	"encoding/json"
	"investigations/db"
	"io"
	"log"
	"net/http"
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

func decode(resp *http.Response, canvas *sql.DB) {

	// close body on function close
	defer resp.Body.Close()

	// make sure request was found
	if resp.StatusCode == 404 {
		return
	}
	log.Printf("%s %s", api.Path, resp.Status)

	// set decoder
	dec := json.NewDecoder(resp.Body)

	// read until break
	for {

		// initalize track info
		var track map[string]interface{}

		// break on EOF
		if err := dec.Decode(&track); err == io.EOF {
			log.Println("JSON", err)
			break
		} else if err != nil {
			panic(err)
		}

		// add track to canvas
		db.AddRow("track", track, canvas)
	}
}

func main() {

	// set the canvas
	canvas := db.InitiateDB()
	defer canvas.Close()

	// start counter at 0
	var trackID int
	for {

		// increment ID
		trackID++

		// attempt to get info on trackID
		api.Path = "tracks/" + strconv.Itoa(trackID)
		if resp, err := http.Get(api.String()); err != nil {
			panic(err)
		} else {

			// decode json
			decode(resp, canvas)
		}
	}
}
