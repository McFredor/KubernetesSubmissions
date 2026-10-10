package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"
)

//go:embed index.html
var indexHTML []byte

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(indexHTML)
	})

	fmt.Printf("Server start in orin %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
