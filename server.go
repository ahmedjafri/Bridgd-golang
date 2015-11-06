package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Listening on server 3000")

	test()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.ListenAndServe(":3000", nil)
}
