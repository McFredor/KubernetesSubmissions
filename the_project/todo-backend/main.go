package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"unicode/utf8"
)

const maxLength = 140

type Todo struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3002"
	}

	var mu sync.Mutex
	todos := []Todo{}
	nextID := 1

	http.HandleFunc("GET /todos", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		list := make([]Todo, len(todos))
		copy(list, todos)
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list)
	})

	http.HandleFunc("POST /todos", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Content string `json:"content"`
		}
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&body); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		content := strings.TrimSpace(body.Content)
		if content == "" || utf8.RuneCountInString(content) > maxLength {
			http.Error(w, fmt.Sprintf("el todo debe tener entre 1 y %d caracteres", maxLength), http.StatusBadRequest)
			return
		}

		mu.Lock()
		todo := Todo{ID: nextID, Content: content}
		nextID++
		todos = append(todos, todo)
		mu.Unlock()

		log.Printf("Todo creado: %q", content)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(todo)
	})

	fmt.Printf("Server started in port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
