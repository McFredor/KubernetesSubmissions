package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	id := rand.Text()

	filePath := os.Getenv("FILE_PATH")
	if filePath == "" {
		filePath = "/usr/src/app/files/log.txt"
	}

	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fmt.Printf("Writer started, writing to %s\n", filePath)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		line := fmt.Sprintf("%s: %s\n", time.Now().UTC().Format(time.RFC3339Nano), id)
		if _, err := f.WriteString(line); err != nil {
			log.Fatal(err)
		}
		fmt.Print(line)
		<-ticker.C
	}
}
