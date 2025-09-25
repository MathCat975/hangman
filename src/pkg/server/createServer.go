package server

import (
	"net/http"
)

func CreateServer() *http.Server {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	return &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}
}

func BootServer() {
	server := CreateServer()
	server.ListenAndServe()
}
