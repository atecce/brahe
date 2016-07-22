package lyricsdotnet

import (
	"database/sql"
	"fmt"
	"investigations/db"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3" // need this to declare sqlite3 pointer
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
)

// set wait group
var wg sync.WaitGroup

// get url
const url = "http://www.lyrics.net"

// constant flags
const href = "href"
const strong = "strong"

// set caught up variable
var caughtUp bool

func communicate(url string) (bool, *html.Node) {

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
		fmt.Println(time.Now(), url, resp.Status)

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

func inASCIIupper(start string) bool {
	for _, char := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		if string(char) == string(start[0]) {
			return true
		}
	}
	return false
}

// Investigate called to start web scrape
func Investigate(start string) {

	// initiate db
	canvas := db.InitiateDB("lyrics_net")

	// use specified start letter
	var expression string
	if start == "0" || start == "" {
		expression = "^/artists/[0A-Z]$"
	} else if inASCIIupper(start) {
		expression = "^/artists/[" + string(start[0]) + "-Z]$"
	} else {
		log.Println("Invalid start character.")
		return
	}

	// set body
	skip, root := communicate(url)

	// check for skip
	if skip {
		return
	}

	letterNodes := scrape.FindAll(root, func(n *html.Node) bool {
		letters, _ := regexp.Compile(expression)
		return letters.MatchString(scrape.Attr(n, "href"))
	})

	// TODO need better iterator name
	for _, n := range letterNodes {

		// concatenate the url TODO almost certainly a better way to join URL's
		letterURL := url + scrape.Attr(n, "href") + "/99999"

		// get artists
		getArtists(start, letterURL, canvas)
	}
}

func getArtists(start, letterURL string, canvas *sql.DB) {

	// set caught up expression
	expression, _ := regexp.Compile("^" + start + ".*$")
	if start == "0" {
		caughtUp = true
	}

	// set body
	skip, root := communicate(letterURL)

	// check for skip
	if skip {
		return
	}

	artistNodes := scrape.FindAll(root, func(n *html.Node) bool {
		artists, _ := regexp.Compile("^artist/.*$")
		if n.Parent != nil {
			return n.Parent.Data == "strong" && artists.MatchString(scrape.Attr(n, "href"))
		}
		return false
	})

	for _, n := range artistNodes {

		// TODO again, must be much better way to join URL's
		artistURL := url + "/" + scrape.Attr(n, "href")

		// artist name
		artistName := scrape.Text(n)

		// check if caught up
		if expression.MatchString(artistName) {
			caughtUp = true
		}
		if !caughtUp {
			continue
		}

		// parse the artist
		parseArtist(artistURL, artistName, canvas)
	}
}

func parseArtist(artistURL, artistName string, canvas *sql.DB) {

	// initialize artist flag
	var artistAdded bool

	// set body
	skip, root := communicate(artistURL)

	// check for skip
	if skip {
		return
	}

	albumNodes := scrape.FindAll(root, func(n *html.Node) bool {
		return scrape.Attr(n, "class") == "artist-album-label"
	})

	for _, n := range albumNodes {

		// TODO awk would be nice here
		text := scrape.Text(n)
		albumTitle := text[:len(text)-7]
		albumYear, _ := strconv.Atoi(text[len(text)-5 : len(text)-1])

		// TODO better urljoin
		albumURL := url + scrape.Attr(n.FirstChild, "href")

		fmt.Println(albumURL, albumTitle, albumYear)

		// add artist
		if !artistAdded {
			db.AddArtist(artistName, canvas)
			artistAdded = true
		}

		// add album
		db.AddAlbum(artistName, albumTitle, canvas)
	}

	// 					// parse album
	// 					dorothy := parsealbum(albumurl, albumtitle, canvas)
	//
	// 					// handle dorothy
	// 					if dorothy {
	// 						noplace(albumtitle, z, canvas)
}

// func noPlace(albumTitle string, z *html.Tokenizer, canvas *sql.DB) {
//
// 	// parse album from artist page
// 	for {
// 		z.Next()
// 		t := z.Token()
// 		switch t.Data {
//
// 		// check for finished album
// 		case "div":
//
// 			for _, a := range t.Attr {
// 				if a.Key == "class" && a.Val == "clearfix" {
// 					wg.Wait()
// 					return
// 				}
// 			}
//
// 		// check for song links
// 		case strong:
//
// 			z.Next()
//
// 			for _, a := range z.Token().Attr {
// 				if a.Key == href {
//
// 					// concatenate the url
// 					songURL := url + a.Val
//
// 					// next token is artist name
// 					z.Next()
// 					songTitle := z.Token().Data
//
// 					// parse song
// 					wg.Add(1)
// 					go parseSong(songURL, songTitle, albumTitle, canvas)
// 				}
// 			}
// 		}
// 	}
// }
//
// func parseAlbum(albumURL, albumTitle string, canvas *sql.DB) bool {
//
// 	// initialize flag that checks for songs
// 	var hasSongs bool
//
// 	// set body
// 	skip, root := communicate(albumURL)
//
// 	// check for skip
// 	if skip {
// 		return false
// 	}
//
// 	// parse page
// 	z := html.NewTokenizer(b)
// 	for {
// 		switch z.Next() {
//
// 		// end of html document
// 		case html.ErrorToken:
// 			wg.Wait()
// 			return !hasSongs
//
// 		// catch start tags
// 		case html.StartTagToken:
//
// 			// check token
// 			t := z.Token()
// 			switch t.Data {
//
// 			// check for home page
// 			case "body":
// 				for _, a := range t.Attr {
// 					if a.Key == "id" && a.Val == "s4-page-homepage" {
// 						return true
// 					}
// 				}
//
// 			// find song links
// 			case strong:
// 				z.Next()
// 				for _, a := range z.Token().Attr {
// 					if a.Key == href {
//
// 						// mark that the page has songs
// 						hasSongs = true
//
// 						// concatenate the url
// 						songURL := url + a.Val
//
// 						// next token is artist name
// 						z.Next()
// 						songTitle := z.Token().Data
//
// 						// parse song
// 						wg.Add(1)
// 						go parseSong(songURL, songTitle, albumTitle, canvas)
// 					}
// 				}
// 			}
// 		}
// 	}
// }
//
// func parseSong(songURL, songTitle, albumTitle string, canvas *sql.DB) {
//
// 	// set body
// 	skip, root := communicate(songURL)
//
// 	// check for skip
// 	if skip {
// 		return
// 	}
//
// 	// parse page
// 	z := html.NewTokenizer(b)
// 	for {
// 		switch z.Next() {
//
// 		// end of html document
// 		case html.ErrorToken:
// 			wg.Done()
// 			return
//
// 		// catch start tags
// 		case html.StartTagToken:
//
// 			// find pre tokens
// 			if z.Token().Data == "pre" {
//
// 				// next token is lyrics
// 				z.Next()
// 				lyrics := z.Token().Data
//
// 				// add song to db
// 				db.AddSong(albumTitle, songTitle, lyrics, canvas)
// 			}
// 		}
// 	}
// }
