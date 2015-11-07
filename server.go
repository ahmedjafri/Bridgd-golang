package main

import (
	"log"
	"net/http"
	"encoding/json"
)

func main() {
	log.Println("Listening on server 3000")

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.HandleFunc("/rooms", func(w http.ResponseWriter, r *http.Request) {
		var rooms []Room = getRooms()

		json.NewEncoder(w).Encode(rooms)
	})

	http.ListenAndServe(":3000", nil)
}
