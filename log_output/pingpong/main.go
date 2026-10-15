package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	filePath := os.Getenv("FILE_PATH")
	if filePath == "" {
		filePath = "/usr/src/app/files/pingpong.txt"
	}

	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		log.Fatal(err)
	}

	// Resume from the count stored in the volume, if any.
	var mu sync.Mutex
	var count int64
	if content, err := os.ReadFile(filePath); err == nil {
		count, _ = strconv.ParseInt(strings.TrimSpace(string(content)), 10, 64)
	}

	http.HandleFunc("GET /pingpong", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		n := count
		count++
		err := os.WriteFile(filePath, []byte(strconv.FormatInt(count, 10)), 0o644)
		mu.Unlock()
		if err != nil {
			log.Printf("failed to save count: %v", err)
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "pong %d\n", n)
	})

	fmt.Printf("Server started in port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
