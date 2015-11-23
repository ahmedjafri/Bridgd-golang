package main

import (
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"io/ioutil"
	"fmt"
	"github.com/jmoiron/sqlx"
	"encoding/json"
)

type SongJSON map[string]interface{}

type Room struct {
	// capitalized variables so they are exported https://golang.org/ref/spec#Exported_identifiers
	Id int `json:"id" db:"id"`
    Name string `json:"name" db:"name"`
    CurrentIndex int `json:"currentIndex" db:"current_index"`
    Queue []SongJSON `json:"queue"`
    Current SongJSON `json:"current"`
}


type Song struct { 
    VideoData string `db:"videoData"`
}

var schema = `
CREATE TABLE rooms (
	id integer not null primary key, 
	name text unique,
	integer current_index
);

CREATE TABLE songs (
	roomId integer not null, 
	videoData text
)`

/*
Method definitions:
room->getSongQueue() []Song
room->DeleteSong(videoLink) 
room->AddSong(videoLink)
*/

// keep a versioning scheme for the db so we can know when to recreate the sqlite file
const dbVersion string = "1.3";

// TODO (ajafri): we perform a file read on each one of these calls so use it sparingly or change the pattern
func createNewDB() (*sqlx.DB, error) {
	var version string = "dbVersion"

	b, err := ioutil.ReadFile("bridgdDBVersion.txt")
    if err != nil {
        b = []byte("") // simulate a version mismatch if reading the file throws an error
    }

    version = string(b)

    if version != dbVersion {
    	log.Println	("Version difference detected - recreating database")

		os.Remove("./rooms.db") // clear the db if it is not versioned to the current version
		
		// let's create a new db and instantiate the structure 
		db, err := sqlx.Open("sqlite3", "./rooms.db") 
    	if err != nil {
    		return db,err
    	}

    	/* TABLE CREATION */
		_, err = db.Exec(schema)
		if err != nil {
			log.Printf("%q: %s\n", err, schema)
			return db, err
		}

		room := Room{Id:1, Name:"no-yolo-zone"}
		// TODO (ajafri): testing code to initialize table with data. Take out. 
		_, err = db.NamedExec("insert into rooms(id, name) values(:id, :name)",room )
		if err != nil {
			log.Printf("%q: %s\n", err)
			return db, err
		}

		// TODO (ajafri): testing code to initialize table with data. Take out. 
		song := Song{VideoData: "{ \"id\": { \"kind\": \"youtube#video\", \"videoId\": \"IDKMKBmpwrg\" }, \"snippet\": { \"title\": \"Tinashe - Player (Audio) ft. Chris Brown\", \"description\": \"\\\"Player\\\" feat. Chris Brown from Tinashe's forthcoming new album, Joyride. Apple Music: http://smarturl.it/PlayerCBa?IQid=yt Spotify: http://smarturl.it/PlayerCBs?\", \"thumbnails\": { \"default\": { \"url\": \"https://i.ytimg.com/vi/IDKMKBmpwrg/default.jpg\" } }, \"channelTitle\": \"TinasheOfficialVEVO\" } }"}
		
		_, err = db.NamedExec("insert into songs(roomId, videoData) values(:id, :videoData)", map[string]interface{}{ "id":room.Id, "videoData":song.VideoData })
		if err != nil {
			log.Printf("%q: %s\n", err)
			return db, err
		}

		err = ioutil.WriteFile("bridgdDBVersion.txt", []byte(dbVersion), 0644)
	    if err != nil {
	        return db, err
	    }
    } 
		
	return sqlx.Open("sqlite3", "./rooms.db")
}

func getRooms() ([]Room) {
	db, err := createNewDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var rooms []Room

	rows, err := db.Queryx("select id, name from rooms")
	defer rows.Close()

    for rows.Next() {
    	var room Room 
        err := rows.StructScan(&room)
        if err != nil {
            log.Fatalln(err)
        } 
        room.Queue = getSongsForRoom(room)
        room.Current = room.Queue[0]
		rooms = append(rooms, room)
    }

	return rooms
}

func getRoom(roomName string) (Room, error) {
	db, err := createNewDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var room Room
    err = db.Get(&room, "select id, name from rooms where name=$1", roomName)

    if(err != nil) {
    	return room, err
    }

	return room, nil
}

func getSongsForRoom(room Room) ([]SongJSON) {
	db, err := createNewDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("select videoData from songs where roomId=%d",room.Id))
	if err != nil {
		log.Fatal(err)
	}

	var songs []SongJSON

	defer rows.Close()
	for rows.Next() {
		var song Song
		rows.Scan(&song.VideoData)

		var songJSON SongJSON
		
		err := json.Unmarshal([]byte(song.VideoData), &songJSON)
		if err != nil {
    		log.Fatal(err)
		}

		songs = append(songs, songJSON)
	}

	return songs
}