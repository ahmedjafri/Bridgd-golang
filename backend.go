package main

import (
	"log"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"fmt"
	"io/ioutil"
)

func serveGetRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	room, err := getRoom(vars["roomName"])

	if(err != nil) {
		fmt.Fprint(w, err)
		return 
	}

	var roomQueueJSON string;
	roomQueueJSONbytes, err := json.Marshal(&room)
	if(err != nil) {
		log.Fatal(err)
	}

	roomQueueJSON = string(roomQueueJSONbytes)
	err = roomTemplate.Execute(w, map[string]string {"Room": roomQueueJSON})
	if(err != nil) {
		log.Fatal(err)
	}
}

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
