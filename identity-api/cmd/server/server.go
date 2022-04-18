package main

import (
	"fmt"
	"net/http"

	"blab.com/library/acl"
)

const Port = ":8070"

func main() {
	fmt.Printf("Starting server on %s\n", Port)
	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(Port, nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	if !acl.HasAccess(r.URL.Path[1:]) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	fmt.Fprintf(w, "Hello Identity API, %s!", r.URL.Path[1:])
}
