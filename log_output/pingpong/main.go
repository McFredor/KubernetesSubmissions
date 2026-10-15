package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	var mu sync.Mutex
	var count int64

	http.HandleFunc("GET /pingpong", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		n := count
		count++
		mu.Unlock()

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "pong %d\n", n)
	})

	// Returns the number of pongs served so far, without incrementing it.
	http.HandleFunc("GET /pings", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		n := count
		mu.Unlock()

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "%d\n", n)
	})

	fmt.Printf("Server started in port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
