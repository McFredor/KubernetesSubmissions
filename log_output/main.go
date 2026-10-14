package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	id := rand.Text()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			fmt.Printf("%s: %s\n", time.Now().UTC().Format(time.RFC3339Nano), id)
			<-ticker.C
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "%s: %s\n", time.Now().UTC().Format(time.RFC3339Nano), id)
	})

	fmt.Printf("Server started in port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
