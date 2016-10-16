package main

import (
	"encoding/json"
	"net/url"
	"os"
	"strconv"
	"sync"

	"github.com/atecce/brahe/heavens"
	"github.com/gocql/gocql"
)

//const trackID = 5151298

const filename = "favorites.txt"

var (
	wg    sync.WaitGroup
	mutex sync.Mutex
)

// pick up where you left off
func findMaxID(session *gocql.Session) int {

	// initialize maximum
	var maxID int

	if err := session.Query(`SELECT MAX(userID) FROM favorites`).
		Scan(&maxID); err != nil {
		panic(err)
	}

	return maxID
}

func main() {

	// check user input
	if len(os.Args) != 2 {
		println("\nUsage: brahe <braching factor>\n")
		os.Exit(1)
	}

	// branching factor
	b, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	// initialize cluster
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "de_nova_stella"
	cluster.ProtoVersion = 4
	session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// observe the heavens
	for id := (findMaxID(session) / 1000) * 1000; ; id++ {

		// look at b stars at a time
		if id%b == 0 {
			wg.Wait()
		}
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// construct method
			methodPath := "users/" + strconv.Itoa(id)
			method := &url.URL{
				Scheme:   "http",
				Host:     "api.soundcloud.com",
				Path:     methodPath,
				RawQuery: "client_id=" + os.Getenv("CLIENT_ID"),
			}

			// get user info
			body := heavens.Observe(method)
			if body == nil {
				return
			}
			var user map[string]interface{}
			json.Unmarshal(body, &user)
			userID := strconv.FormatFloat(user["id"].(float64), 'f', -1, 64)

			// get favoriters info
			method.Path = methodPath + "/favorites"
			body = heavens.Observe(method)
			var songs []interface{}
			json.Unmarshal(body, &songs)
			for _, song := range songs {
				trackID := strconv.FormatFloat(song.(map[string]interface{})["id"].(float64), 'f', -1, 64)

				// record the observation
				if err := session.Query(`INSERT INTO favorites (userID, trackID) VALUES (?, ?)`,
					userID, trackID).Exec(); err != nil {
					panic(err)
				}
			}
		}(id)
	}
}
