package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "todo-server")
	})

	fmt.Printf("Server start in orin %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
