//
// I should not like my writing to spare other people the trouble of thinking.
// But, if possible, to stimulate someone to thoughts of their own.
//

package lyrics_net

import (
	"db"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// get url
var url string = "http://www.lyrics.net"

func communicate(url string) io.ReadCloser {

	// never stop trying
	for {

		// get url
		resp, err := http.Get(url)

		// catch error
		if err != nil {
			fmt.Println("ERROR: Failed to crawl \"" + url + "\"")
			return nil
		}

		// check status codes
		if resp.StatusCode == 200 {
			return resp.Body
		} else if resp.StatusCode == 403 {
			fmt.Println(url, "Forbidden.")
			return nil
		} else if resp.StatusCode == 404 {
			fmt.Println(url, "Not found.")
			return nil
		} else if resp.StatusCode == 503 {
			fmt.Println(url, "Overloaded server.")
			time.Sleep(30 * time.Minute)
		} else if resp.StatusCode == 504 {
			fmt.Println(url, "Gateway timeout.")
			time.Sleep(time.Minute)
		} else {
			fmt.Println(url, "Other status code:", resp.StatusCode)
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

func Investigate(start string) {

	// initiate db
	db.InitiateDB("lyrics_net")

	// use specified start letter
	var expression string
	if start == "0" || start == "" {
		expression = "^/artists/[0A-Z]$"
	} else if inASCIIupper(start) {
		expression = "^/artists/[" + string(start[0]) + "-Z]$"
	} else {
		fmt.Println("Invalid start character.")
		return
	}

	// set regular expression for letter suburls
	letters, _ := regexp.Compile(expression)

	// set body
	b := communicate(url)
	defer b.Close()

	// declare tokenizer
	z := html.NewTokenizer(b)

	for {
		// get next token
		next := z.Next()

		switch {

		// catch error
		case next == html.ErrorToken:
			return

		// catch start tags
		case next == html.StartTagToken:

			// set token
			t := z.Token()

			// find a tokens
			if t.Data == "a" {

				// iterate over token
				for _, a := range t.Attr {

					// if the link is inside
					if a.Key == "href" {

						// and the link matches the letters
						if letters.MatchString(a.Val) {

							// concatenate the url
							letter_url := url + a.Val + "/99999"

							// get artists
							getArtists(start, letter_url)
						}
					}
				}
			}
		}
	}
}

func getArtists(start, letter_url string) {

	// set caught up expression
	expression, _ := regexp.Compile("^" + url + "/artist/" + start + ".*/[0-9]*$")
	var caught_up bool
	if start == "0" {
		caught_up = true
	}

	// set regular expression for letter suburls
	artists, _ := regexp.Compile("^artist/.*$")

	// set body
	b := communicate(letter_url)
	defer b.Close()

	// declare tokenizer
	z := html.NewTokenizer(b)

	for {

		// get next token
		next := z.Next()

		switch {

		// catch error
		case next == html.ErrorToken:
			return

		// catch start tags
		case next == html.StartTagToken:

			// find strong tokens
			if z.Token().Data == "strong" {

				// get next token
				z.Next()

				// iterate over token
				for _, a := range z.Token().Attr {

					// if the link is inside
					if a.Key == "href" {

						// if it matches the artist string
						if artists.MatchString(a.Val) {

							// concatenate the url
							artist_url := url + "/" + a.Val

							// check if caught up
							if expression.MatchString(artist_url) {
								caught_up = true
							}
							if !caught_up {
								continue
							}

							// next token is artist name
							z.Next()
							artist_name := z.Token().Data

							// parse the artist
							parseArtist(artist_url, artist_name)
						}
					}
				}
			}
		}
	}
}

func parseArtist(artist_url, artist_name string) {

	// initialize artist flag
	var artistAdded bool

	// set body
	b := communicate(artist_url)
	fmt.Println()
	fmt.Println(artist_name, artist_url)
	defer b.Close()

	// declare tokenizer
	z := html.NewTokenizer(b)

	for {
		// get next token
		next := z.Next()

		switch {

		// catch error
		case next == html.ErrorToken:
			return

		// catch start tags
		case next == html.StartTagToken:

			// set token
			t := z.Token()

			// look for artist album labels
			if t.Data == "h3" {
				for _, a := range t.Attr {
					if a.Key == "class" && a.Val == "artist-album-label" {

						// add artist
						if !artistAdded {
							db.AddArtist(artist_name)
							artistAdded = true
						}

						// album links are next token
						var album_url string
						z.Next()
						for _, album_attribute := range z.Token().Attr {
							if album_attribute.Key == "href" {
								album_url = url + album_attribute.Val
							}
						}

						// album titles are the next token
						z.Next()
						album_title := z.Token().Data

						// add album
						db.AddAlbum(artist_name, album_title)

						// parse album
						dorothy := parseAlbum(album_url, album_title)

						// handle dorothy
						if dorothy {

							// set flag for finished album
							var finished bool

							for {
								// set next token
								z.Next()
								t = z.Token()

								// check for finished album
								if t.Data == "div" {
									for _, a := range t.Attr {
										if a.Key == "class" && a.Val == "clearfix" {
											finished = true
										}
									}
								}
								if finished {
									break
								}

								// check for song link
								if t.Data == "strong" {
									z.Next()
									for _, a := range z.Token().Attr {
										if a.Key == "href" {

											// concatenate the url
											song_url := url + a.Val

											// next token is artist name
											z.Next()
											song_title := z.Token().Data

											// parse song
											parseSong(song_url, song_title, album_title)
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func parseAlbum(album_url, album_title string) bool {

	// set body
	b := communicate(album_url)
	fmt.Println()
	fmt.Println("\t", album_title, album_url)
	defer b.Close()

	// declare tokenizer
	z := html.NewTokenizer(b)

	for {
		// get next token
		next := z.Next()

		switch {

		// catch error
		case next == html.ErrorToken:
			return false

		// catch start tags
		case next == html.StartTagToken:

			// set token
			t := z.Token()

			// check for home page
			if t.Data == "body" {
				for _, a := range t.Attr {
					if a.Key == "id" && a.Val =="s4-page-homepage" {
						return true
					}
				}
			}

			// find strong tokens
			if t.Data == "strong" {

				// get next token
				z.Next()
				t = z.Token()

				// iterate over token
				for _, a := range t.Attr {

					// if the link is inside
					if a.Key == "href" {

						// concatenate the url
						song_url := url + a.Val

						// next token is artist name
						z.Next()
						song_title := z.Token().Data

						// parse song
						parseSong(song_url, song_title, album_title)
					}
				}
			}
		}
	}
}

func parseSong(song_url, song_title, album_title string) {

	// set body
	b := communicate(song_url)
	fmt.Println()
	fmt.Println("\t\t\t", song_title, song_url)
	defer b.Close()

	// declare tokenizer
	z := html.NewTokenizer(b)

	for {
		// get next token
		next := z.Next()

		switch {

		// catch error
		case next == html.ErrorToken:
			return

		// catch start tags
		case next == html.StartTagToken:

			// find pre tokens
			if z.Token().Data == "pre" {

				// next token is lyrics
				z.Next()
				lyrics := z.Token().Data

				// print lyrics
				fmt.Println()
				for _, line := range strings.Split(lyrics, "\n") {
					fmt.Println("\t\t\t\t", line)
				}
				fmt.Println()

				// add song to db
				db.AddSong(album_title, song_title, lyrics)
			}
		}
	}
}
