package main

import (
	"fmt"
	"net/http"
)

func NewServer() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Face Recognition Server Running")
	})
	return mux
}
