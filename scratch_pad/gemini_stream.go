package scratch_pad

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"log"
	"os"
	"strings"
)

func GeminiStreamRun() {
	// Load Gemini API key from environment variable
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		fmt.Println("GEMINI_API_KEY environment variable not set")
		return
	}

	// Set API endpoint and context
	ctx := context.Background()

	// Create a Generative Language API client
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal("Error creating Generative Language API client:", err)
		return
	}
	defer func(client *genai.Client) {
		err := client.Close()
		if err != nil {

		}
	}(client)

	model := client.GenerativeModel("gemini-1.5-flash")

	fmt.Print("Enter your message: ")
	reader := bufio.NewReader(os.Stdin)
	query, _ := reader.ReadString('\n')
	query = strings.TrimSpace(query)

	//// Simulate loading indicator
	//for i := 0; i < 3; i++ {
	//	fmt.Print(".")
	//	time.Sleep(time.Second / 2)
	//}
	//fmt.Println()
	//
	//// Count tokens
	//queryTokens, err := client.GenerativeModel("gemini-1.5-flash").CountTokens(ctx, genai.Text(query))
	//if err != nil {
	//	fmt.Println("Error counting tokens:", err)
	//	return
	//}

	iter := model.GenerateContentStream(ctx, genai.Text(query))

	// Print request and total token information
	//fmt.Printf("Your Request (**%d tokens**): %s\n", queryTokens.TotalTokens, query)

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
			// Check if candidate has text content
			if candidate.Content != nil && candidate.Content.Parts != nil {
				for _, part := range candidate.Content.Parts {
					fmt.Printf("%v", part)
					//totalResponseTokens += len(part)
				}
			}
		}
	}
	//res, err = cs.SendMessage(ctx, genai.Text("Tell me a joke in 300 words"))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(res)
}
