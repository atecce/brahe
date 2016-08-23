package main

import (
	"bodhi/db"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

//const trackID = 5151298

// initalize api url
var api = &url.URL{
	Scheme:   "http",
	Host:     "api.soundcloud.com",
	RawQuery: "client_id=" + os.Getenv("CLIENT_ID"),
}

func communicate(api *url.URL, canvas *sql.DB) {

	// never stop trying
	for {

		if resp, err := http.Get(api.String()); err != nil {
			panic(err)
		} else {
			defer resp.Body.Close()
			log.Printf("%s %s", api.Path, resp.Status)

			switch resp.StatusCode {
			case 404:
				return
			case 502:
				time.Sleep(time.Minute)
			default:
				decode(resp, canvas)
				return
			}
		}
	}
}

func decode(resp *http.Response, canvas *sql.DB) {

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
	canvas := db.Initiate()
	defer canvas.Close()

	// start counter at 0
	var trackID int
	db.GetLatest(&trackID, canvas)
	for {

		// increment ID
		trackID++

		// attempt to get info on trackID
		api.Path = "tracks/" + strconv.Itoa(trackID)

		// try and communicate
		communicate(api, canvas)
	}
}
