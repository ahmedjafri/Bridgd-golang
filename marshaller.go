package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"io/ioutil"
)

type Room struct {
	// capitalized variables so they are exported https://golang.org/ref/spec#Exported_identifiers
	Id int `json:"id"`
    Name string `json:"name"`
}


type Song struct {
    videoData string // do we need this?
}

/*
Method definitions:
room->getSongQueue() []Song
room->DeleteSong(videoLink) 
room->AddSong(videoLink)
*/

// keep a versioning scheme for the db so we can know when to recreate the sqlite file
const dbVersion string = "1";

// TODO (ajafri): we perform a file read on each one of these calls so use it sparingly or change the pattern
func createNewDB() (*sql.DB, error) {
	var version string = "dbVersion"

	b, err := ioutil.ReadFile("bridgdDBVersion.txt")
    if err != nil {
        b = []byte("") // simulate a version mismatch if reading the file throws an error
    }

    version = string(b)

    if version != dbVersion {
    	log.Println	("Version difference detected - recreating database")

	    err = ioutil.WriteFile("bridgdDBVersion.txt", []byte(dbVersion), 0644)
	    if err != nil {
	        panic(err)
	    }

		os.Remove("./rooms.db") // clear the db if it is not versioned to the current version
		
		// let's create a new db and instantiate the structure 
		db, err := sql.Open("sqlite3", "./rooms.db") 
    	if err != nil {
    		return db,err
    	}

    	var query string = "create table rooms (id integer not null primary key, name text);"
		_, err = db.Exec(query)
		if err != nil {
			log.Printf("%q: %s\n", err, query)
			return db, err
		}

		query = "insert into rooms(id, name) values(1, \"no yolo zone\")"
		_, err = db.Exec(query)
		if err != nil {
			log.Printf("%q: %s\n", err, query)
			return db, err
		}

    } 
		
	return sql.Open("sqlite3", "./rooms.db")
}

func getRooms() ([]Room) {
	db, err := createNewDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("select id, name from rooms")
	if err != nil {
		log.Fatal(err)
	}

	var rooms []Room

	defer rows.Close()
	for rows.Next() {
		var room Room
		rows.Scan(&room.Id, &room.Name)
		rooms = append(rooms, room)
	}

	return rooms
}