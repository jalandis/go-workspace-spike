package main

import (
	"fmt"
	"net/http"
)

const Port = ":8071"

func main() {
	fmt.Printf("Starting server on %s\n", Port)
	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(Port, nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Assessment API, %s!", r.URL.Path[1:])
}
