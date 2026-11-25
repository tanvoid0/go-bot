package server

import (
	"fmt"
	"github.com/tanvoid0/go_bot/bot"
	"github.com/tanvoid0/go_bot/data"
	"io"
	"net/http"
	"strconv"
)

// Assuming NewGeminiBot() and StreamResponse() are already implemented
// geminiBot := bot.NewGeminiBot()
var (
	client = bot.GetClient()
)

func BotStreamHandler(r *http.Request) (io.Reader, error) {
	// Read the prompt from the request body
	prompt := r.FormValue("prompt")
	if prompt == "" {
		return nil, fmt.Errorf("missing prompt in request body")
	}
	roomId := r.FormValue("roomId")
	var room *data.ChatRoom
	var err error // Declare the error variable outside to avoid shadowing

	if roomId == "" {
		// Create a new chat room if roomId is empty
		room = &data.ChatRoom{}
		_, err = room.Create() // Use the already declared err variable
		if err != nil {
			return nil, err // Return error if creating room fails
		}
	} else {
		// Parse roomId as int
		id, parseErr := strconv.Atoi(roomId) // Parse the string into an int
		if parseErr != nil {
			return nil, fmt.Errorf("invalid roomId format: %v", parseErr)
		}

		// Find the existing chat room by ID
		room, err = data.FindChatRoomById(id) // Use the parsed int for the FindChatRoomById method
		if err != nil {
			return nil, err // Return error if room is not found
		}
	}
	return client.AskStream(prompt)
}
