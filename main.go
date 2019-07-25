package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

var DepotRoot = os.Getenv("DEPOT_ROOT")

// Check that DEPOT_ROOT is defined and can be accessed
func checkRoot() {
	if DepotRoot == "" {
		log.Fatalf("DEPOT_ROOT not defined")
	}

	info, err := os.Stat(DepotRoot)
	if err != nil {
		log.Fatalf("DEPOT_ROOT error: %s", err)
	} else if !info.IsDir() {
		log.Fatalf("DEPOT_ROOT is not a directory")
	}
}

// Server setup
func Server() {
	// Check setup
	checkRoot()

	// Routing
	r := mux.NewRouter()
	r.HandleFunc("/", RootHandler)
	r.HandleFunc("/jobs/{job}/{filepath:.*}", JobsHandler)
	r.HandleFunc("/zip", ZipHandler)

	// Port
	addr := ":8080"
	DepotPort := os.Getenv("DEPOT_PORT")
	if DepotPort != "" {
		addr = fmt.Sprintf(":%s", DepotPort)
	}

	// Serve
	log.Printf("running at http://127.0.0.1%s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

func main() {
	Server()
}
