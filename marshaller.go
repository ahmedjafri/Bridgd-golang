package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

type Room struct {
	id int
    name string
}


type Song struct {
    name string // do we need this?
    videoLink string // Can we retrieve all the information from this link? Do we need the other params
    videoData string // do we need this?
}

/*
Method definitions:
room->getSongQueue() []Song
room->DeleteSong(videoLink) 
room->AddSong(videoLink)
*/


func createNewDB() (*sql.DB, error) {
	// TODO (ajafri): Should we be clearing the database on startup. Probably not.
	os.Remove("./rooms.db")
	return sql.Open("sqlite3", "./rooms.db")
}


// TODO (ajafri): This is just an example of how we can read/write to the db. Take it out. 
func writeToDB() {
	db, err := createNewDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table rooms (id integer not null primary key, name text);
	delete from rooms;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("insert into rooms(id, name) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// underscore here means discard the resultant. A variable you can write to but not read from. It will get compiled out.
	_, err = stmt.Exec(1, fmt.Sprintf("This is bridgd room %d!!", 1))
	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()

	rows, err := db.Query("select id, name from rooms")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		fmt.Println(id, name)
	}
}