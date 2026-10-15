package main

import (
	_ "embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

//go:embed index.html
var indexHTML []byte

const (
	imageURL = "https://picsum.photos/1200"
	cacheTTL = 10 * time.Minute
)

var (
	imagePath   string
	mu          sync.Mutex
	downloading bool
)

// downloadImage descarga una imagen nueva y la guarda de forma atómica en imagePath.
// La fecha de modificación del archivo marca cuándo se obtuvo la imagen.
func downloadImage() error {
	resp, err := http.Get(imageURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status inesperado: %s", resp.Status)
	}

	tmp, err := os.CreateTemp(filepath.Dir(imagePath), "image-*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	if _, err := io.Copy(tmp, resp.Body); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	log.Println("Imagen nueva descargada")
	return os.Rename(tmp.Name(), imagePath)
}

// refreshInBackground descarga una imagen nueva sin bloquear la petición actual.
func refreshInBackground() {
	mu.Lock()
	if downloading {
		mu.Unlock()
		return
	}
	downloading = true
	mu.Unlock()

	go func() {
		defer func() {
			mu.Lock()
			downloading = false
			mu.Unlock()
		}()
		if err := downloadImage(); err != nil {
			log.Println("Error descargando imagen:", err)
		}
	}()
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	info, err := os.Stat(imagePath)
	if os.IsNotExist(err) {
		// No hay imagen en caché: hay que descargarla antes de responder.
		if err := downloadImage(); err != nil {
			http.Error(w, "no se pudo obtener la imagen", http.StatusBadGateway)
			log.Println("Error descargando imagen:", err)
			return
		}
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if time.Since(info.ModTime()) > cacheTTL {
		// Imagen vencida: se sirve la vieja una última vez y se pide la nueva.
		refreshInBackground()
	}

	w.Header().Set("Cache-Control", "no-store")
	http.ServeFile(w, r, imagePath)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	imagePath = os.Getenv("IMAGE_PATH")
	if imagePath == "" {
		imagePath = "/usr/src/app/files/image.jpg"
	}
	if err := os.MkdirAll(filepath.Dir(imagePath), 0o755); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(indexHTML)
	})
	http.HandleFunc("/image", imageHandler)

	// Test Persistence.
	http.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Apagando el contenedor...")
		log.Println("Apagado solicitado vía /shutdown")
		go func() {
			time.Sleep(100 * time.Millisecond)
			os.Exit(1)
		}()
	})

	fmt.Printf("Server start in port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
