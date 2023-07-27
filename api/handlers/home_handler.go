package handlers

import "net/http"

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world 2"))
}

func ExampleHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello worldsdcasdcadcadca"))
}
