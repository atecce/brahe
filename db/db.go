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
		log.Println("Failed to open db", name + ".db:", err)
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
		log.Println("Failed to create tables:", err)
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
		log.Println("Failed to add artist", artist_name + ":", err)
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
		log.Println("Failed to add album", album_title, "by", artist_name+":", err)
	}
}

func AddSong(album_title, song_title, lyrics string) {

	for {
		var failed bool

		// prepare db
		db := PrepareDB()
		defer db.Close()
		tx, err := db.Begin()

		// catch error
		if err != nil {
			failed = true
			log.Println("\ndb.Begin failed.\nFailed to add song", song_title, "in album", album_title+":", err, "\n")
			db.Close()
			time.Sleep(time.Second)
			continue
		}

		// prepare statement
		stmt, err := tx.Prepare("insert or replace into songs (album_title, title, lyrics) values (?, ?, ?)")
		defer stmt.Close()

		// catch error
		if err != nil {
			failed = true
			log.Println("\ntx.Prepare failed.\nFailed to add song", song_title, "in album", album_title+":", err, "\n")
			db.Close()
			time.Sleep(time.Second)
			continue
		}

		// execute statement
		_, err = stmt.Exec(album_title, song_title, lyrics)
		tx.Commit()

		// catch error
		if err != nil {
			failed = true
			log.Println("\nstmt.Exec failed.\nFailed to add song", song_title, "in album", album_title+":", err, "\n")
			db.Close()
			time.Sleep(time.Second)
			continue
		}

		// notify that a previous failure was cleaned up
		if failed {
			log.Println("Successfully added song", song_title, "in album", album_title, err)
		}

		// exit
		break
	}
}
