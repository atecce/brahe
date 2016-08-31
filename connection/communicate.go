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
	Canvas *db.Canvas
}

func (api *API) Communicate(table string, method *url.URL) {

	// never stop trying
	for {

		resp, err := http.Get(method.String())
		if err != nil {
			log.Println(err.Error())
			time.Sleep(time.Minute)
			continue
		}
		defer resp.Body.Close()
		log.Printf("%s %s", method.Path, resp.Status)

		switch resp.StatusCode {
		case 404:
			api.Canvas.AddMissing(method.Path)
			return
		case 500:
			time.Sleep(time.Minute)
		case 502:
			time.Sleep(time.Minute)
		default:
			api.decode(table, resp.Body)
			return
		}
	}
}

func (api *API) decode(table string, body io.Reader) {

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
			time.Sleep(time.Minute)
			panic(err)
		}

		// add track to canvas
		api.Canvas.AddRow(table, track)
	}
}
