package main

import (
	"log"
	"net/http"
	"encoding/json"
	"github.com/googollee/go-socket.io"
	"text/template"
	"github.com/gorilla/mux"
)

func main() {
	log.Println("Initializing socket.io server...")

	// use the default transports by passing in nil
	server, err := socketio.NewServer(nil)
    if err != nil {
        log.Fatal(err)
    }

    server.On("joined", func(so socketio.Socket) {
        log.Println("Connected to a peer...")

        so.On("disconnection", func() {
            log.Println("Disconnected from peer...")
        })
    })   

    r := mux.NewRouter()
    r.Handle("/socket.io/", server)

	createNewDB()
	
	fs := http.FileServer(http.Dir("static/assets/"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))
	
	r.HandleFunc("/rooms", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		   case "GET":
		      	var rooms []Room = getRooms()
				json.NewEncoder(w).Encode(rooms)
		   default:
		   		var error string = "No method at this endpoint" 
		   		json.NewEncoder(w).Encode(error)
		}
	})


    t, err := template.ParseFiles("static/index.html")
    if err != nil {
        log.Fatal(err)
    }

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { 
		var room Room = getRooms()[0]
		var roomQueueJSON string;
		roomQueueJSONbytes, err := json.Marshal(&room)
		if(err != nil) {
			log.Fatal(err)
		}

		roomQueueJSON = string(roomQueueJSONbytes)
		err = t.Execute(w, map[string]string {"Room": roomQueueJSON})
		if(err != nil) {
			log.Fatal(err)
		}


	})

	r.HandleFunc("/emit", func(w http.ResponseWriter, r *http.Request) {
		server.BroadcastTo("chat", "chat message")
	})



	var port string = "3001" // needs to be string so we can concatenate it without converting
	
	log.Println("Listening on server " + port)
	http.Handle("/", r)
	http.ListenAndServe(":" + port, r)
}
