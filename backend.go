package main

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"fmt"
	"io/ioutil"
)

// add a song onto the queue
func servePostRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	room, err := getRoom(vars["roomName"])

	if err != nil {
		fmt.Fprint(w, err)
		return 
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	room, err = addSongToRoom(string(body), room)

	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	json.NewEncoder(w).Encode(room)
}
