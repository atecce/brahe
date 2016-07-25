package websites

import (
	"database/sql"
	"fmt"
	"investigations/connection"
	"investigations/db"
	"log"
	"regexp"
	"strconv"
	"sync"

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
	skip, root := communication.Communicate(url)

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
	skip, root := communication.Communicate(letterURL)

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
	skip, root := communication.Communicate(artistURL)

	// check for skip
	if skip {
		return
	}

	albumNodes := scrape.FindAll(root, func(n *html.Node) bool {
		return scrape.Attr(n, "class") == "artist-album-label"
	})

	for _, n := range albumNodes {

		// album link is first child
		albumTitle := scrape.Text(n.FirstChild)
		albumURL := url + scrape.Attr(n.FirstChild, "href") // TODO better urljoin

		// album year is last child
		albumYearText := scrape.Text(n.LastChild)
		albumYear, _ := strconv.Atoi(albumYearText[1 : len(albumYearText)-1])

		fmt.Println(albumTitle, albumYear, albumURL)

		// add artist
		if !artistAdded {
			db.AddArtist(artistName, canvas)
			artistAdded = true
		}

		// add album
		db.AddAlbum(artistName, albumTitle, albumYear, canvas)

		// parse album
		dorothy := parseAlbum(albumURL, albumTitle, canvas)

		// handle dorothy
		if dorothy {
			noPlace(albumTitle, n, canvas)
		}
	}
}

func noPlace(albumTitle string, titleNode *html.Node, canvas *sql.DB) {

	// get album root node from title node
	albumRoot, _ := scrape.FindParent(titleNode, func(n *html.Node) bool {
		return scrape.Attr(n, "class") == "clearfix"
	})

	// get table nodes
	tableNodes := scrape.FindAll(albumRoot, func(n *html.Node) bool {
		return n.Data == "tr"
	})

	// fields are a slice of strings
	var fields []string

	for i, n := range tableNodes {

		// first node contains the field titles
		if i == 0 {

			// extract all the field nodes
			fieldNodes := scrape.FindAll(n, func(n *html.Node) bool {
				return n.Data == "th"
			})

			for _, fieldNode := range fieldNodes {

				// extract the field
				field := scrape.Text(fieldNode)

				// add non-empty fields
				if field != "" {
					fields = append(fields, field)
				}
			}
		} else {

			// extract all the song nodes
			songNodes := scrape.FindAll(n, func(n *html.Node) bool {
				return scrape.Attr(n, "class") == "tal qx"
			})

			songData := make(map[string]string)

			for i, songNode := range songNodes {

				// set song url and title
				if fields[i] == "Song" && songNode.FirstChild.Data == "strong" {
					songTitle := scrape.Text(songNode)
					songURL := url + scrape.Attr(songNode.FirstChild.FirstChild, "href")
					fmt.Println(songTitle, songURL)
				}

				// set song data
				songData[fields[i]] = scrape.Text(songNode)
			}
		}
	}
}

// 			for _, a := range t.Attr {
// 				if a.Key == "class" && a.Val == "clearfix" {
// 					wg.Wait()
// 					return
// 				}
// 			}
//
// 					// parse song
// 					wg.Add(1)
// 					go parseSong(songURL, songTitle, albumTitle, canvas)
// 				}
// 			}
// 		}
// 	}
// }

func parseAlbum(albumURL, albumTitle string, canvas *sql.DB) bool {

	// initialize flag that checks for songs
	var hasSongs bool

	// set body
	skip, root := communication.Communicate(albumURL)

	// check for homepage
	_, skip = scrape.Find(root, func(n *html.Node) bool {
		return scrape.Attr(n, "id") == "s4-page-homepage"
	})

	// check for skip
	if skip {
		return true
	}

	// extract song links with bold
	songNodes := scrape.FindAll(root, func(n *html.Node) bool {
		if n.Parent != nil && scrape.Attr(n, "href") != "" {
			return n.Parent.Data == "strong"
		}
		return false
	})

	for _, n := range songNodes {

		hasSongs = true

		// set song title
		songTitle := scrape.Text(n)

		// set song url
		songURL := url + scrape.Attr(n, "href")

		fmt.Println(songTitle, songURL)
	}

	return !hasSongs

	// // parse page
	// z := html.NewTokenizer(b)
	// for {
	// 	switch z.Next() {
	//
	// 	// end of html document
	// 	case html.ErrorToken:
	// 		wg.Wait()
	// 		return !hasSongs
	//
	// 	// catch start tags
	// 	case html.StartTagToken:
	//
	// 		// check token
	// 		t := z.Token()
	// 		switch t.Data {
	//
	// 		// check for home page
	// 		case "body":
	// 			for _, a := range t.Attr {
	// 				if a.Key == "id" && a.Val == "s4-page-homepage" {
	// 					return true
	// 				}
	// 			}
	//
	// 		// find song links
	// 		case strong:
	// 			z.Next()
	// 			for _, a := range z.Token().Attr {
	// 				if a.Key == href {
	//
	// 					// mark that the page has songs
	// 					hasSongs = true
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
}

//
// func parseSong(songURL, songTitle, albumTitle string, canvas *sql.DB) {
//
// 	// set body
// 	skip, root := communication.Communicate(songURL)
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
