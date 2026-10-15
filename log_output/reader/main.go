package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	filePath := os.Getenv("FILE_PATH")
	if filePath == "" {
		filePath = "/usr/src/app/files/log.txt"
	}

	pingpongPath := os.Getenv("PINGPONG_FILE_PATH")
	if pingpongPath == "" {
		pingpongPath = "/usr/src/app/files/pingpong.txt"
	}

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		content, err := os.ReadFile(filePath)
		if err != nil {
			http.Error(w, "log file not available yet", http.StatusServiceUnavailable)
			return
		}
		lines := strings.Split(strings.TrimSpace(string(content)), "\n")
		lastLine := lines[len(lines)-1]

		pings := "0"
		if c, err := os.ReadFile(pingpongPath); err == nil {
			pings = strings.TrimSpace(string(c))
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "%s.\nPing / Pongs: %s\n", lastLine, pings)
	})

	fmt.Printf("Server started in port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
