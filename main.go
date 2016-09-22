package main

import (
	"encoding/json"
	"net/url"
	"os"
	"strconv"

	"github.com/atecce/brahe/canvas"
	"github.com/atecce/brahe/heavens"
)

//const trackID = 5151298

const (
	table  = "users"
	family = "favorites"
)

var deNovaStella = &canvas.Canvas{}

func main() {

	// set the canvas
	deNovaStella.Initiate()
	defer deNovaStella.Close()
	deNovaStella.AddTable(table)
	deNovaStella.AddFamily(table, family)

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

		// get track info
		body := heavens.Observe(method)
		if body == nil {
			continue
		}
		var user map[string]interface{}
		json.Unmarshal(body, &user)
		row := user["permalink"].(string)

		// get favoriters info
		method.Path = methodPath + "/" + family
		body = heavens.Observe(method)
		var elements []interface{}
		json.Unmarshal(body, &elements)
		for _, element := range elements {
			column := element.(map[string]interface{})["permalink"].(string)
			deNovaStella.Record(table, row, family, column)
		}
	}
}
