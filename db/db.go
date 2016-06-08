package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

var name string

func PrepareDB() *sql.DB {

	// create db
	db, err := sql.Open("sqlite3", name + ".db")

	// catch error
	if err != nil {
		log.Println("ERROR: Failed to open db:", err)
	}

	return db
}

func InitiateDB(db_name string) {

	name = db_name

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
		log.Println("ERROR: Failed to create tables:", err)
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
		log.Println("ERROR: Failed to add artist:", err)
		log.Println(artist_name)
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
		log.Println("ERROR: Failed to add album:", err)
		log.Println(artist_name, album_title)
	}
}

func AddSong(album_title, song_title, lyrics string) {

	for {

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
			log.Println("ERROR: Failed to add song:", err)
			log.Println(album_title, song_title)
			time.Sleep(time.Second)
			continue
		}

		break
	}
}
