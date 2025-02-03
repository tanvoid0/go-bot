package bot

import (
	"bufio"
	"fmt"
	"github.com/tanvoid0/go-bot/util"
	"io"
	"log"
	"net/http"
)

// RunBotCli processes the bot CLI arguments and executes the corresponding commands
func RunBotCli(args []string) {
	if args[0] != "bot" {
		fmt.Println("Usage: dev_pilot bot [model] [command] [options]")
		return
	}
}

// BotModel defines the interface for AI tools
type BotModel interface {
	SetProperties(properties BotProperties)
	Respond(prompt string) (string, error)
	StreamRespond(prompt string) (io.ReadCloser, error)
}

// BotProperties holds configuration parameters for the bot
type BotProperties struct {
	apiKey      string
	model       string
	temperature float32
	tone        string
	instruction string
}

func GetBotProperties() *BotProperties {
	return &BotProperties{
		temperature: 1,
		tone:        "",
		instruction: "A personal help bot",
	}
}

func RunBot() {
	for {
		prompt := util.Read("Enter your message: ")
		if prompt == "exit" {
			break
		}
		geminiBot := NewGeminiBot()

		stream, err := geminiBot.StreamResponse(prompt)
		if err != nil {
			fmt.Println(err)
			continue
		}
		scanner := bufio.NewScanner(stream)
		//var response string // Accumulate the streamed text

		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		//fmt.Println(util.MarkdownFormat(response))
	}
}

func BotHttpStreamRequest(w http.ResponseWriter, r *http.Request) {
	prompt := r.FormValue("prompt")
	if prompt == "" {
		http.Error(w, "Missing prompt in request body", http.StatusBadRequest)
		return
	}

	geminiBot := NewGeminiBot()
	stream, err := geminiBot.StreamResponse(prompt)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	scanner := bufio.NewScanner(stream)

	for scanner.Scan() {
		_, err := fmt.Fprintf(w, "%s\n", scanner.Text())
		if err != nil {
			http.Error(w, "Failed to send event", http.StatusInternalServerError)
			return
		}

		w.(http.Flusher).Flush()
	}

	if err := scanner.Err(); err != nil {
		http.Error(w, "Error reading stream", http.StatusInternalServerError)
	}
}
