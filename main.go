package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// test this with
// curl -N https://minimal-sse-iy4vzwh2ta-ew.a.run.app

func main() {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			panic("Streaming not supported!")
		}

		fmt.Fprintf(w, "data: %s\n\n", time.Now().Format(time.RFC3339))

		ticker := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-r.Context().Done():
				fmt.Println("SSE done")
				return
			case <-ticker.C:
				fmt.Fprintf(w, "data: %s\n\n", time.Now().Format(time.RFC3339))
				fmt.Println("SSE sent data")
				flusher.Flush()
			}
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server := &http.Server{
		Addr:    ":" + port,
		Handler: h2c.NewHandler(r, &http2.Server{}),
	}
	log.Println("Serving http://localhost:" + port)
	log.Fatal(server.ListenAndServe())
}
