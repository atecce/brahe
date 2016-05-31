//
// I should not like my writing to spare other people the trouble of thinking.
// But, if possible, to stimulate someone to thoughts of their own.
//

package main

import (
	"io"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"regexp"
)

// get url
var url string = "http://www.lyrics.net"

func getArtists(letter_url string) {

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

							// next token is artist name
							z.Next()
							artist_name := z.Token()

							// display
							fmt.Println()
							fmt.Println(artist_name, artist_url)
							fmt.Println()

							// parse the artist
							parseArtist(artist_url)

						}
					}
				}
			}
		}
	}
}

func parseArtist(artist_url string) {

	// set body
	b := communicate(artist_url)
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

						// album links are next token
						z.Next()
						for _, album_attribute := range z.Token().Attr {
							if album_attribute.Key == "href" {
								fmt.Println("\t", url + album_attribute.Val)
							}
						}

						// album titles are the next token
						z.Next()
						fmt.Println("\t", z.Token())
					}
				}
			}
		}
	}
}

func communicate(url string) io.ReadCloser {

	// get url
	resp, err := http.Get(url)

	// catch error
	if err != nil {
		fmt.Println("ERROR: Failed to crawl \"" + url + "\"")
		return nil
	}

	// return body
	return resp.Body
}

func main() {

	// set regular expression for letter suburls
	letters, _ := regexp.Compile("^/artists/[0A-Z]$")

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
							getArtists(letter_url)

						}
					}
				}
			}
		}
	}
}

//class lyrics_site:	
//
//    # keep track of everything
//    processes = list()
//
//    siblings = {"artists": list(),
//                "albums":  list(),
//                "songs":   list()}
//
//    def __init__(self, start, branch, verbose):
//
//        # specify flags
//        self.start   = start
//        self.branch  = branch
//        self.verbose = verbose
//
//    def multitask(self, level, function, process_name, process_args):
//
//        # fork the process
//        process = Process(target=function, name=process_name, args=process_args)
//        process.start()
//
//        # track the processes
//        self.processes.append(process)
//        self.siblings[level].append(process)
//
//        print self.branch
//
//        # pace yourself
//        while len(self.siblings[level]) >= self.branch:
//
//            print "Pruning", level, "..."
//            self.siblings[level].pop(0).join()
//
//class lyrics_net(lyrics_site): 
//
//    # know where you are
//    url = 'http://www.lyrics.net/'
//
//    canvas = canvas("lyrics_net")
//
//    def investigate(self):
//
//        # get the soup
//        soup = self.communicate(self.url)
//
//        # set artist expression
//        expression = str()
//
//        if self.start == '0': expression = '^/artists/[0A-Z]$'
//
//        elif self.start[0] in string.ascii_uppercase:
//
//            expression = '^/artists/[' + self.start[0] + '-Z]$'
//
//        # pick up where you left off
//        caught_up = bool()
//
//        # for each artist
//        for artist_name, artist_url in artist_data: 
//
//            # check if you've caught up
//            if re.match("^" + self.url + "artist/" + self.start + ".*/[0-9]*$", artist_url): caught_up = True
//
//            # if you haven't, continue
//            if not caught_up: continue
//
//            # fork
//            self.multitask('artists', self.honor, artist_name, (artist_name, artist_url,))
//
//    def honor(self, artist_name, artist_url):
//
//        # get the soup
//        artist_soup = self.communicate(artist_url)
//
//        # add the artist to the canvas
//        self.canvas.add_artist(artist_name)
//
//        # get the album items
//        album_items = artist_soup.find_all('div', {'class': 'clearfix'})
//
//        # for each item
//        for item in album_items: 
//
//            # check for a header
//            if item.h3:
//
//                # check the header contains a link
//                if item.h3.a:
//
//                    # set the album information
//                    album_title = item.h3.a.text
//                    album_url   = urljoin(self.url, item.h3.a.get('href'))
//
//                    if self.verbose:
//
//                        print '\t', album_title
//                        print '\t', album_url
//                        print
//
//                    # add the album to the canvas
//                    self.canvas.add_album(artist_name, album_title)
//
//                    # get the soup
//                    album_soup = self.communicate(album_url)
//
//                    # handle Dorothy (which do not return the proper status code)
//                    if album_soup.find_all('body', {'id': 's4-page-homepage'}): 
//
//                        # extract the song data
//                        song_data = ((trace.a.text, urljoin(self.url, trace.a.get('href'))) \
//                                      for trace in item.find_all('tr') 	   		    \
//                                                if trace.a)
//
//                    # otherwise
//                    else:
//
//                        # extract the song data
//                        song_data = ((song_tag.a.text, urljoin(self.url, song_tag.a.get('href'))) \
//                                      for song_tag in album_soup.find_all('strong') 	   	  \
//                                                   if song_tag.a)
//
//                    # for each song
//                    for song_title, song_url in song_data:
//
//                        # fork
//                        self.multitask('songs', self.meditate, song_title, (album_title, song_title, song_url,))
//
//    def meditate(self, album_title, song_title, song_url):
//
//        # make some soup
//        song_soup = self.communicate(song_url)
//
//        # sometimes there's nothing to meditate on
//        try: 
//
//            lyrics = song_soup.find_all('pre', {'id': 'lyric-body-text'})[0].text
//
//            if self.verbose:
//
//                for line in lyrics.splitlines(): print '\t\t\t', line
//
//                print
//
//        except IndexError: return
//
//        # add song to canvas
//        self.canvas.add_song(album_title, song_title, lyrics)
//
//if __name__ == '__main__':
//
//    # declare parser
//    parser  = argparse.ArgumentParser()
//
//    # add arguments
//    start   = parser.add_argument("-s", "--start",   help="specify the start character",  default='0')
//    branch  = parser.add_argument("-b", "--branch",  help="specify the branching factor", default=2, 
//                                 type=int)
//    verbose = parser.add_argument("-v", "--verbose", help="specify the verbosity",        default=False,
//                                 action='store_true')
//
//    # parse arguments
//    args = parser.parse_args()
//
//    # start the investigation
//    investigation = lyrics_net(args.start, args.branch, args.verbose)
//    investigation.investigate()
//
//    # shut down instance after finished
//    system("sudo shutdown -h now")
