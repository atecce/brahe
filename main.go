package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

//const trackID = 5151298

func main() {

	// start counter at 0
	var trackID int
	for {

		// get json
		if resp, err := http.Get("http://api.soundcloud.com/tracks/" + strconv.Itoa(trackID) +
			"?client_id=" + os.Getenv("CLIENT_ID")); err != nil {
			panic(err)
		} else {
			defer resp.Body.Close()

			// increment ID
			fmt.Println(trackID)
			trackID++

			// set decoder
			dec := json.NewDecoder(resp.Body)

			// read until break
			for {

				// initalize track info
				var track map[string]interface{}

				// break on err (hopefully EOF)
				if err := dec.Decode(&track); err != nil {
					log.Println(err)
					break
				}

				// iterate through keys and print values
				for k := range track {
					fmt.Println(k, track[k])
				}
			}
		}
	}
}
