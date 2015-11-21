package main

import (
	"log"
	"net/http"
	"encoding/json"
)

func main() {
	createNewDB()
	
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.HandleFunc("/rooms", func(w http.ResponseWriter, r *http.Request) {
		var rooms []Room = getRooms()
		json.NewEncoder(w).Encode(rooms)
	})

	var port string = "3001" // needs to be string so we can concatenate it without converting
	
	log.Println("Listening on server " + port)
	http.ListenAndServe(":" + port, nil)
}
