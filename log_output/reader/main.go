package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
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

	pingpongURL := os.Getenv("PINGPONG_URL")
	if pingpongURL == "" {
		pingpongURL = "http://pingpong-svc:2346/pings"
	}

	client := &http.Client{Timeout: 2 * time.Second}

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		content, err := os.ReadFile(filePath)
		if err != nil {
			http.Error(w, "log file not available yet", http.StatusServiceUnavailable)
			return
		}
		lines := strings.Split(strings.TrimSpace(string(content)), "\n")
		lastLine := lines[len(lines)-1]

		pings, err := fetchPings(client, pingpongURL)
		if err != nil {
			log.Printf("failed to fetch pings: %v", err)
			http.Error(w, "pingpong service not available", http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "%s.\nPing / Pongs: %s\n", lastLine, pings)
	})

	fmt.Printf("Server started in port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func fetchPings(client *http.Client, url string) (string, error) {
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(body)), nil
}
