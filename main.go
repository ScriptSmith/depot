package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"path"
)

// Check that DEPOT_ROOT is defined and can be accessed
func checkRoot(root string) {
	if root == "" {
		log.Fatalf("DEPOT_ROOT not defined")
	}

	info, err := os.Stat(root)
	if err != nil {
		log.Fatalf("DEPOT_ROOT error: %s", err)
	} else if !info.IsDir() {
		log.Fatalf("DEPOT_ROOT is not a directory")
	}

	tmpDir := path.Join(root, "tmp")
	info, err = os.Stat(tmpDir);
	if err != nil {
		err = os.Mkdir(tmpDir, 0777)
		if err != nil {
			log.Fatalf("Couldn't create tmp dir")
		}
	} else if !info.IsDir() {
		log.Fatalf("tmp dir is not a directory")
	}
}

// Get the mux router
func getRouter(handlers *Handlers) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.RootHandler)
	r.HandleFunc("/assets/{filepath:.*}", handlers.AssetsHandler)
	r.HandleFunc("/jobs/{job}/{filepath:.*}", handlers.JobsHandler)
	r.HandleFunc("/zip", handlers.ZipHandler)
	return r
}

// Server setup
func Server() {
	// Check setup
	root := os.Getenv("DEPOT_ROOT")
	checkRoot(root)
	user := os.Getenv("DEPOT_USER")
	pass := os.Getenv("DEPOT_PASS")
	deletion := os.Getenv("DEPOT_DISABLE_DELETION")

	// Create handlers with root
	handlers := &(Handlers{
		root:     root,
		deletion: deletion == "",
	})

	// Define auth
	auth := &(Auth{
		user: user,
		pass: pass,
	})

	// Routing
	r := getRouter(handlers)

	// Auth
	r.Use(auth.Middleware)

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
