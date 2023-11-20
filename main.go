package main

import (
	"log"
)

func main() {
	// Create a new PostgreSQL store by calling NewPostgresStore.
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	// run ther server on port:3000
	server := newAPIServer(":3000", store)
	server.Run()
}
