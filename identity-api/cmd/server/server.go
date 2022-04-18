package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"blab.com/library/acl"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8060"
		log.Printf("defaulting to port %v\n", port)
	}

	fmt.Printf("Starting identity server on %s\n", port)
	http.HandleFunc("/", IdentityServer)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

func IdentityServer(w http.ResponseWriter, r *http.Request) {
	if !acl.HasAccess(r.URL.Path[1:]) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	fmt.Fprintf(w, "Hello Identity API, %s!", r.URL.Path[1:])
}
