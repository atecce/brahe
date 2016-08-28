package connection

import (
	"bodhi/db"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type API struct {
	Method *url.URL
	Canvas *db.Canvas
}

func (api *API) Communicate() {

	// never stop trying
	for {

		resp, err := http.Get(api.Method.String())
		if err != nil {
			log.Println(err.Error())
			panic(err)
		}
		defer resp.Body.Close()
		log.Printf("%s %s", api.Method.Path, resp.Status)

		switch resp.StatusCode {
		case 404:
			api.Canvas.AddMissing(api.Method.Path)
			return
		case 502:
			time.Sleep(time.Minute)
		default:
			api.decode(resp.Body)
			return
		}
	}
}

func (api API) decode(body io.Reader) {

	// set decoder
	dec := json.NewDecoder(body)

	// read until break
	for {

		// initalize track info
		var track map[string]interface{}

		// break on EOF
		if err := dec.Decode(&track); err == io.EOF {
			log.Println(err)
			break
		} else if err != nil {
			log.Println(err.Error())
			panic(err)
		}

		// add track to canvas
		api.Canvas.AddRow("track", track)
	}
}
