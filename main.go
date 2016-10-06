package main

import (
	"encoding/json"
	"log"
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

// var deNovaStella = &canvas.Canvas{}

func openFile(filename string) *os.File {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	return f
}

func main() {

	// set the canvas
	// deNovaStella.Initiate()
	// defer deNovaStella.Close()
	// deNovaStella.AddTable(table)
	// deNovaStella.AddFamily(table, family)

	favorites := openFile("favorites.txt")
	defer favorites.Close()

	id, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	// observe the heavens
	for ; ; id++ {

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
			log.Println("trackID:", trackID)
			_, err := favorites.WriteString(userID + "\t" + trackID + "\n")
			if err != nil {
				panic(err)
			}
		}
	}
}
