package main

import (
	"log"
	"net/http"
	"encoding/json"
	"text/template"
	htemplate "html/template" 
	"github.com/gorilla/mux"
	"fmt"
	"io/ioutil"
)

var roomTemplate *template.Template
var indexTemplate *htemplate.Template
var err error

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

func serveRoom(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "GET": 
			if(roomTemplate == nil) {
				roomTemplate, err = template.ParseFiles("static/room.html")
				if err != nil {
				    log.Fatal(err)
				}
			}
			serveGetRoom(w,r)

		case "POST": 
			servePostRoom(w,r)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprint(w,"the method " + r.Method + " is not supported")
	}



}

func serveHome(w http.ResponseWriter, r *http.Request) { 
	if(indexTemplate == nil) {
		indexTemplate, err = htemplate.ParseFiles("static/index.html")
		if err != nil {
		    log.Fatal(err)
		}
	}

	var rooms []Room = getRooms()

	err = indexTemplate.Execute(w, map[string][]Room {"Rooms": rooms})
	if(err != nil) {
		log.Fatal(err)
	}
}