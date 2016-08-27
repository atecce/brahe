package connection

import (
	"bodhi/db"
	"bodhi/herodotus"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

var httpLog = herodotus.CreateFileLog("http")

func Communicate(api *url.URL, canvas *sql.DB) {

	// never stop trying
	for {

		resp, err := http.Get(api.String())
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		httpLog.Printf("%s %s", api.Path, resp.Status)

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

func decode(resp *http.Response, canvas *sql.DB) {

	// set decoder
	dec := json.NewDecoder(resp.Body)

	// read until break
	for {

		// initalize track info
		var track map[string]interface{}

		// break on EOF
		if err := dec.Decode(&track); err == io.EOF {
			httpLog.Println(err)
			break
		} else if err != nil {
			panic(err)
		}

		// add track to canvas
		db.AddRow("track", track, canvas)
	}
}
