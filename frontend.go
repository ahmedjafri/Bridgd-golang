package main

import (
	"log"
	"net/http"
	"text/template"
	htemplate "html/template" 
	"fmt"
)

var roomTemplate *template.Template
var indexTemplate *htemplate.Template
var err error

func serveRoom(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "GET": 
			if(roomTemplate == nil) {
				roomTemplate, err = template.ParseFiles("templates/room.html")
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
		indexTemplate, err = htemplate.ParseFiles("templates/index.html")
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