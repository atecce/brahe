package connection

import (
	"bodhi/db"
	"bodhi/herodotus"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

type API struct {
	Method *url.URL
	Canvas *db.Canvas
}

var httpLog = herodotus.CreateFileLog("http")

func (api *API) Communicate() {

	// never stop trying
	for {

		resp, err := http.Get(api.Method.String())
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		httpLog.Printf("%s %s", api.Method.Path, resp.Status)

		switch resp.StatusCode {
		case 404:
			return
		case 502:
			time.Sleep(time.Minute)
		default:
			api.decode(resp)
			return
		}
	}
}

func (api API) decode(resp *http.Response) {

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
		api.Canvas.AddRow("track", track)
	}
}
