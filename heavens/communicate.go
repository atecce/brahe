package heavens

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

func Observe(method *url.URL) []byte {

	// never stop trying
	for {

		resp, err := http.Get(method.String())
		if err != nil {
			log.Println(err) // TODO
			time.Sleep(time.Minute)
			continue
		}
		defer resp.Body.Close()
		fmt.Printf("%s %s\n", method.Path, resp.Status) // TODO

		switch resp.StatusCode {
		case 403, 404: // TODO
			return nil
		case 500, 502, 503:
			time.Sleep(time.Minute)
		default:
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println(err) // TODO
				time.Sleep(time.Minute)
				continue
			}
			return body
		}
	}
}
