package connection

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/atecce/brahe/db"
)

type API struct {
	Canvas *db.Canvas
}

func (api *API) Communicate(table string, method *url.URL) {

	// never stop trying
	for {

		resp, err := http.Get(method.String())
		if err != nil {
			// TODO need more intelligent logging
			log.Println(err.Error())
			time.Sleep(time.Minute)
			continue
		}
		defer resp.Body.Close()
		// TODO need more intelligent logging
		log.Printf("%s %s", method.Path, resp.Status)

		switch resp.StatusCode {
		case 403:
			return
		case 404:
			// api.Canvas.AddMissing(method.Path)
			return
		case 500:
			time.Sleep(time.Minute)
		case 502:
			time.Sleep(time.Minute)
		default:
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}
			api.decode(table, body)
			return
		}
	}
}

func (api *API) decode(table string, body []byte) {

	// marshal json to type check
	var row map[string]interface{}
	err := json.Unmarshal(body, &row)
	if err != nil {
		panic(err)
	}
	// TODO need more intelligent logging
	log.Println(row)

	// add track to canvas
	api.Canvas.AddRow(table, row)
}
