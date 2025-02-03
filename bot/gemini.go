package bot

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"github.com/tanvoid0/go-bot/data"
	"github.com/tanvoid0/go-bot/util"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"io"
	"log"
	"reflect"
)

type GeminiBot struct {
	properties BotProperties
	ctx        context.Context
	client     genai.Client
	model      genai.GenerativeModel
}

func NewGeminiBot() *GeminiBot {
	properties := GetBotProperties()
	properties.apiKey = util.ReadEnv("GEMINI_API_KEY")
	properties.model = "gemini-1.5-flash"

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(properties.apiKey))
	if err != nil {
		log.Fatalln("Error creating Generative Language API client:", err)
		return nil
	}
	model := client.GenerativeModel(properties.model)

	return &GeminiBot{
		properties: *properties,
		ctx:        context.Background(),
		client:     *client,
		model:      *model,
	}
}

func (bot GeminiBot) StreamResponse(prompt string) (io.Reader, error) {

	user := data.User{
		Name:     "Tanveer",
		Username: "tanvoid0",
		Password: "password",
	}
	message := data.ChatMessage{
		Data:   prompt,
		Sender: user.ID,
	}
	_, err := message.Create()
	if err != nil {
		return nil, err
	}
	iter := bot.model.GenerateContentStream(bot.ctx, genai.Text(prompt))
	buffer := make(chan string) // create a channel for streaming string

	go func() {
		defer close(buffer) // close the channel after processing is done
		for {
			resp, err := iter.Next()
			if errors.Is(err, iterator.Done) {
				break
			}
			if err != nil {
				fmt.Println("Error during streaming:", err)
				log.Fatal(err)
				return
			}
			for _, candidate := range resp.Candidates {
				if candidate.Content != nil && candidate.Content.Parts != nil {
					for _, part := range candidate.Content.Parts {
						// Handle different content types based on the concrete type of `part`
						switch p := part.(type) {
						case genai.Text:
							buffer <- string(p)
						case genai.Blob:
							// Handle Blob data (e.g., image, audio)
							fmt.Println("Received Blob:", p.MIMEType)
						case genai.FunctionCall, genai.FunctionResponse, genai.ExecutableCode, genai.CodeExecutionResult:
							// Handle other content types as needed (if applicable)
							fmt.Println("Unknown content type:", reflect.TypeOf(p))
						default:
							// Handle unexpected types (optional)
							fmt.Println("Unexpected Part type:", reflect.TypeOf(p))
						}
					}
				}
			}
		}
	}()
	return readerFromChannel(buffer), nil // return a reader from the channel
}

func readerFromChannel(ch chan string) io.Reader {
	return &channelReader{ch: ch} // create a custom reader from the channel
}

type channelReader struct {
	ch <-chan string
}

func (r *channelReader) Read(p []byte) (n int, err error) {
	text, ok := <-r.ch
	if !ok {
		return 0, io.EOF
	}
	return copy(p, text), nil
}
