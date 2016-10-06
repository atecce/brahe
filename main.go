package main

import (
	"encoding/json"
	"net/url"
	"os"
	"strconv"

	"github.com/atecce/brahe/heavens"
)

//const trackID = 5151298

const (
	table  = "users"
	family = "favorites"
)

func main() {

	// initalize file
	favorites, err := os.Create("favorites.txt")
	if err != nil {
		panic(err)
	}
	defer favorites.Close()

	// observe the heavens
	for id := 1; ; id++ {

		// construct method
		methodPath := table + "/" + strconv.Itoa(id)
		method := &url.URL{
			Scheme:   "http",
			Host:     "api.soundcloud.com",
			Path:     methodPath,
			RawQuery: "client_id=" + os.Getenv("CLIENT_ID"),
		}

		// get user info
		body := heavens.Observe(method)
		if body == nil {
			continue
		}
		var user map[string]interface{}
		json.Unmarshal(body, &user)
		userID := strconv.FormatFloat(user["id"].(float64), 'f', -1, 64)

		// get favoriters info
		method.Path = methodPath + "/" + family
		body = heavens.Observe(method)
		var songs []interface{}
		json.Unmarshal(body, &songs)
		for _, song := range songs {
			trackID := strconv.FormatFloat(song.(map[string]interface{})["id"].(float64), 'f', -1, 64)

			// record the observation
			_, err := favorites.WriteString(userID + "\t" + trackID + "\n")
			if err != nil {
				panic(err)
			}
		}
	}
}
