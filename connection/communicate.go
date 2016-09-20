package connection

import (
	"encoding/json"
	"fmt"
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

func (api *API) Communicate(family string, method *url.URL) bool {

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
		fmt.Printf("%s %s\n", method.Path, resp.Status)

		switch resp.StatusCode {
		case 403:
			return true
		case 404:
			// api.Canvas.AddMissing(method.Path)
			return true
		case 500, 502, 503:
			time.Sleep(time.Minute)
		default:
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				// TODO need more intelligent logging
				log.Println(err.Error())
				time.Sleep(time.Minute)
			}
			api.decode(body)

			// add track to canvas
			// api.Canvas.AddEntry(method.Path, family, entry)
			return false
		}
	}
}

func (api *API) decode(body []byte) {

	// marshal json to type check
	var entry interface{}
	err := json.Unmarshal(body, &entry)
	if err != nil {
		panic(err)
	}

	printEntry("", entry)
}

func printEntry(leadingTabs string, entry interface{}) {

	// switch on entry type
	switch entry.(type) {

	// iterate through keys and values on map
	case map[string]interface{}:
		for k, v := range entry.(map[string]interface{}) {

			// switch on value type
			switch v.(type) {

			// recursively call
			case map[string]interface{}, []interface{}:
				fmt.Println(k)
				printEntry(leadingTabs+"\t", v)

			// print keys and values
			default:
				fmt.Println(leadingTabs, k, v)
			}
		}

	// iterate through elements on slice
	case []interface{}:
		for _, element := range entry.([]interface{}) {

			// switch on element type
			switch element.(type) {

			// recursively call
			case map[string]interface{}, []interface{}:
				printEntry(leadingTabs, element)

			// print element of list
			default:
				fmt.Println(leadingTabs, element)
			}
		}
	default:
		panic(entry)
	}
}
