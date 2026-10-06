package main

import (
	"crypto/rand"
	"fmt"
	"time"
)

func main() {
	id := rand.Text()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		fmt.Printf("%s: %s\n", time.Now().UTC().Format(time.RFC3339Nano), id)
		<-ticker.C
	}
}
