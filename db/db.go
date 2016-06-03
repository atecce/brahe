//
// I should not like my writing to spare other people the trouble of thinking.
// But, if possible, to stimulate someone to thoughts of their own.
//

package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func PrepareDB() *sql.DB {

	// create db
	db, err := sql.Open("sqlite3", "lyrics_net.db")

	// catch error
	if err != nil {
		fmt.Println("ERROR: Failed to open db:", err)
	}

	return db
}

func InitiateDB() {

	// prepare db
	db := PrepareDB()
	defer db.Close()

	// create tables
	_, err := db.Exec(`create table if not exists artists (
								  
				  name text not null, 	  
				          			  
				  primary key (name))`)

	_, err = db.Exec(`create table if not exists albums ( 		
									
				 title       text not null, 			
				 artist_name text not null,  		
				       					
				 primary key (title, artist_name), 				
				 foreign key (artist_name) references artists (name))`)

	_, err = db.Exec(`create table if not exists songs ( 	    	       
								       
				 title       text not null, 	    	       
				 album_title text not null, 	    	       
				 lyrics      text, 			    	       
				       				       
				 primary key (album_title, title),
				 foreign key (album_title) references albums (title))`)

	// catch error
	if err != nil {
		fmt.Println("ERROR: Failed to create tables:", err)
	}
}

func AddArtist(artist_name string) {

	// prepare db
	db := PrepareDB()
	defer db.Close()
	tx, err := db.Begin()

	// insert entry
	stmt, err := tx.Prepare("insert or replace into artists (name) values (?)")
	defer stmt.Close()
	_, err = stmt.Exec(artist_name)
	tx.Commit()

	// catch error
	if err != nil {
		fmt.Println("ERROR: Failed to add artist:", err)
	}
}

func AddAlbum(artist_name, album_title string) {

	// prepare db
	db := PrepareDB()
	defer db.Close()
	tx, err := db.Begin()

	// insert entry
	stmt, err := tx.Prepare("insert or replace into albums (artist_name, title) values (?, ?)")
	defer stmt.Close()
	_, err = stmt.Exec(artist_name, album_title)
	tx.Commit()

	// catch error
	if err != nil {
		fmt.Println("ERROR: Failed to add album:", err)
	}
}

func AddSong(album_title, song_title, lyrics string) {

	// prepare db
	db := PrepareDB()
	defer db.Close()
	tx, err := db.Begin()

	// insert entry
	stmt, err := tx.Prepare("insert or replace into songs (album_title, title, lyrics) values (?, ?, ?)")
	defer stmt.Close()
	_, err = stmt.Exec(album_title, song_title, lyrics)
	tx.Commit()

	// catch error
	if err != nil {
		fmt.Println("ERROR: Failed to add song:", err)
	}
}
