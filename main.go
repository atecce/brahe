package main

import (
	"bufio"
	"encoding/json"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/atecce/brahe/heavens"
)

//const trackID = 5151298

const filename = "favorites.txt"

var wg sync.WaitGroup

// pick up where you left off
func findMaxID(favorites *os.File) int {

	// initialize maximum
	var maxID int

	// iterate through file line by line
	scanner := bufio.NewScanner(favorites)
	for scanner.Scan() {

		// extract ID from line
		id, err := strconv.Atoi(strings.Split(scanner.Text(), "\t")[0])
		if err != nil {
			panic(err)
		}

		// test maximum ID
		if id > maxID {
			maxID = id
		}
	}

	// check for scanning errors
	if err := scanner.Err(); err != nil {
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

	// initialize file
	favorites, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		panic(err)
	}
	defer favorites.Close()

	// observe the heavens
	for id := findMaxID(favorites); ; id++ {

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
				_, err := favorites.WriteString(userID + "\t" + trackID + "\n")
				if err != nil {
					panic(err)
				}
			}
		}(id)
	}
}
