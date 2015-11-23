package main

import (
	"log"
	"net/http"
	"github.com/googollee/go-socket.io"
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

	// these routes can be found in frontend.go
	r.HandleFunc("/room/{roomName}", serveRoom)
	r.HandleFunc("/", serveHome)


	// TODO (ajafri): testing. Take this out later
	r.HandleFunc("/emit", func(w http.ResponseWriter, r *http.Request) {
		server.BroadcastTo("chat", "chat message")
	})

	var port string = "3001" // needs to be string so we can concatenate it without converting
	
	log.Println("Listening on server " + port)
	http.Handle("/", r)
	http.ListenAndServe(":" + port, r)
}
