package communication

import (
	"log"
	"net/http"
	"time"

	"golang.org/x/net/html"
)

func Communicate(url string) (bool, *html.Node) {

	// never stop trying
	for {

		// get url
		resp, err := http.Get(url)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second)
			continue
		}
		defer resp.Body.Close()

		// get root node
		root, err := html.Parse(resp.Body)
		if err != nil {
			panic(err)
		}

		// write status to output
		//fmt.Println(time.Now(), url, resp.Status)

		// check status codes
		switch resp.StatusCode {

		// cases which are returned
		case 200:
			return false, root
		case 403:
			return true, root
		case 404:
			return true, root

		// cases which are retried
		case 503:
			time.Sleep(30 * time.Minute)
		case 504:
			time.Sleep(time.Minute)
		default:
			time.Sleep(time.Minute)
		}
	}
}
