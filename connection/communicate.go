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

func (api *API) Communicate(family string, method *url.URL) {

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
				// TODO need more intelligent logging
				log.Println(err.Error())
				time.Sleep(time.Minute)
			}
			entry := api.decode(family, body)

			// add track to canvas
			api.Canvas.AddEntry(method.Path, family, entry)
			return
		}
	}
}

func (api *API) decode(family string, body []byte) map[string]interface{} {

	// marshal json to type check
	var entry map[string]interface{}
	err := json.Unmarshal(body, &entry)
	if err != nil {
		panic(err)
	}

	return entry
}
