package main

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/atecce/brahe/heavens"
	"github.com/gocql/gocql"
)

//const trackID = 5151298

// pick up where you left off
func findMaxID(session *gocql.Session) int {
	var maxID int
	if err := session.Query(`SELECT MAX(userID) FROM favorites`).
		Scan(&maxID); err != nil {
		panic(err)
	}
	return maxID
}

func main() {

	var wg sync.WaitGroup

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
	cluster.ProtoVersion = 4
	cluster.Keyspace = "de_nova_stella"
	session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}

	// listen for interrupt to clean up
	listener := make(chan os.Signal, 1)
	signal.Notify(listener, os.Interrupt)
	go func() {
		<-listener
		log.Println("INFO closing session")
		session.Close()
		os.Exit(0)
	}()

	// figure out where to start
	var startID int
	if maxID := findMaxID(session); maxID < 1000 {
		startID = 0
	} else {
		startID = maxID - 1000
	}

	// observe the heavens
	for id := startID; ; id++ {

		// look at b stars at a time
		if id%b == 0 {
			log.Println("INFO Waiting...")
			wg.Wait()
		}
		wg.Add(1)
		go func(id int, wg *sync.WaitGroup) {
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

				for {

					// record the observation
					if err := session.Query(`INSERT INTO favorites (userID, trackID) VALUES (?, ?)`,
						userID, trackID).Exec(); err != nil {

						if err == gocql.ErrTimeoutNoResponse {
							log.Println("ERROR", err)
							time.Sleep(time.Minute)
							continue
						}

						panic(err)
					}

					break
				}
			}
		}(id, &wg)
	}
}
