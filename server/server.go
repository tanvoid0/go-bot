package server

import (
	"bufio"
	"fmt"
	"github.com/tanvoid0/dev-bot/bot"
	"github.com/tanvoid0/dev-bot/util"
	"io"
	"log"
	"net/http"
)

func streamRequestHandler(stream io.Reader, err error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err != nil {
			http.Error(w, "Failed to get stream response", http.StatusInternalServerError)
			return
		}
		// Setting headers for streaming response
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		scanner := bufio.NewScanner(stream)

		// Stream scanner.Text() to the client
		for scanner.Scan() {
			_, err := fmt.Fprintf(w, "%s ", scanner.Text())
			if err != nil {
				http.Error(w, "Failed to send data to client", http.StatusInternalServerError)
				return
			}
			// Flush the buffer to send the data immediately
			w.(http.Flusher).Flush()
		}

		// Check for scanner errors
		if err := scanner.Err(); err != nil {
			http.Error(w, "Error reading stream", http.StatusInternalServerError)
		}
	}
}

func botStreamHandler(r *http.Request) (io.Reader, error) {
	// Read the prompt from the request body
	prompt := r.FormValue("prompt")
	if prompt == "" {
		return nil, fmt.Errorf("missing prompt in request body")
	}

	// Assuming NewGeminiBot() and StreamResponse() are already implemented
	geminiBot := bot.NewGeminiBot()
	return geminiBot.StreamResponse(prompt)
}

func handlerFunction(w http.ResponseWriter, r *http.Request) {
	stream, err := botStreamHandler(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	streamRequestHandler(stream, nil)(w, r)
}

func Run() {
	// Usage:  curl -X POST -d "prompt=Tell me a story in 1000 words" http://localhost:8000/api/bot/query
	http.HandleFunc("/api/bot/query", handlerFunction) // Use http.HandleFunc here
	port := util.ReadEnv("PORT")

	fmt.Println("Server started at http://localhost:" + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server error:", err)
	}
}
