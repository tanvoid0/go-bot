package bot

import (
	"bufio"
	"fmt"
	"github.com/tanvoid0/go_bot/util"
	"io"
	"log"
	"strings"
	"time"
)

type Client interface {
	GetProperties() *Properties
	Ask(query string) (string, error)
	AskStream(query string) (io.ReadCloser, error)
}

// GetHistory returns the chat history from the bot
func GetHistory(b Client) []Message {
	return b.GetProperties().History
}

// AddToChatHistory adds a message to the bot's chat history
func AddToChatHistory(b Client, msg Message) {
	props := b.GetProperties()
	props.History = append(props.History, msg)
}

type Request struct {
	Model    string    `json:"Model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

// Response struct represents the response from the chat Model
type Response struct {
	Model   string  `json:"Model"`
	Message Message `json:"message"`
	Done    bool    `json:"done"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Properties holds configuration parameters for the bot
type Properties struct {
	client      string
	apiKey      string
	Model       string
	temperature float32
	tone        string
	instruction string
	user        string
	bot         string
	History     []Message
}

func GetBotProperties() *Properties {
	return &Properties{
		temperature: 1,
		tone:        "",
		instruction: "A personal help bot",
		client:      util.ReadEnv("CLIENT"),
		user:        util.ReadEnv("USER"),
		bot:         util.ReadEnv("BOT"),
		Model:       util.ReadEnv("MODEL"),
	}
}

func TrimResponse(data string) string {
	return strings.Join(strings.Fields(data), " ")
}

func RunBot(c Client) {
	for {
		prompt := util.Read("Enter your message: ")
		if prompt == "exit" {
			break
		}
		//geminiBot := NewGeminiBot()

		start := time.Now()
		fmt.Println("Started queries")

		stream, err := c.AskStream(prompt)
		if err != nil {
			fmt.Println(err)
			continue
		}
		scanner := bufio.NewScanner(stream)
		//var response string // Accumulate the streamed text

		for scanner.Scan() {
			fmt.Printf("%s ", scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("\nCompleted in %v\n\n", time.Since(start))

		//fmt.Println(util.MarkdownFormat(response))
	}
}

func GetClient() Client {
	var client Client = nil
	// Check the CLIENT environment variable
	if util.ReadEnv("CLIENT") == "ollama" {
		client = NewOllamaClient() // Return OllamaClient
	} else {
		// Handle the case where CLIENT is not recognized, for example, logging and returning an error
		log.Fatal("Invalid client type specified in environment variable 'CLIENT'")
	}

	// Safely retrieve properties from the client
	if client == nil {
		log.Fatal("Client is nil. Exiting.")
	}

	properties := client.GetProperties()
	log.Printf("Running Bot: %s, Model: %s\n", properties.bot, properties.client)

	return client
}

//// RunBotCli processes the bot CLI arguments and executes the corresponding commands
//func RunBotCli(args []string) {
//	if args[0] != "bot" {
//		fmt.Println("Usage: dev_pilot bot [Model] [command] [options]")
//		return
//	}
//}
