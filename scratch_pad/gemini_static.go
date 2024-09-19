package scratch_pad

import (
	"bufio"
	"context"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"log"
	"os"
	"strings"
)

func GeminiRun(args []string) {
	ctx := context.Background()

	apiKey, ok := os.LookupEnv("GEMINI_API_KEY")
	if !ok {
		log.Fatalln("Environment variable GEMINI_API_KEY not set")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	model.SetTemperature(1)
	model.SetTopK(64)
	model.SetTopP(0.95)
	model.SetMaxOutputTokens(8192)
	model.ResponseMIMEType = "text/plain"

	// model.SafetySettings = Adjust safety settings
	// See https://ai.google.dev/gemini-api/docs/safety-settings

	session := model.StartChat()
	session.History = []*genai.Content{}

	fmt.Print("Enter your message: ")
	reader := bufio.NewReader(os.Stdin)
	command, _ := reader.ReadString('\n')
	command = strings.TrimSpace(command)

	resp, err := session.SendMessage(ctx, genai.Text(command))
	if err != nil {
		log.Fatalf("Error sending message: %v", err)
	}

	for _, part := range resp.Candidates[0].Content.Parts {
		fmt.Printf("%v\n", part)
	}
}

// initialize the chat
//cs := model.StartChat()
//cs.History = []*genai.Content{
//	&genai.Content{
//		Parts: []genai.Part{
//			genai.Text("Hello, I am glad to meet you"),
//		},
//		Role: "user",
//	},
//	&genai.Content{
//		Parts: []genai.Part{
//			genai.Text("It's nice to meet you too. What's on your mind?"),
//		},
//		Role: "model",
//	},
//}
