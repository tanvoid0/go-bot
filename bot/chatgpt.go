package bot

import (
	"io"
	"time"
)

type ChatGPTBot struct{}

func (c *ChatGPTBot) Response(prompt string) (string, error) {
	return "ChatGPT response to: " + prompt, nil
}

func (c *ChatGPTBot) StreamResponse(prompt string) (io.ReadCloser, error) {
	reader, writer := io.Pipe()
	go func() {
		defer writer.Close()
		for i := 0; i < 5; i++ {
			writer.Write([]byte("ChatGPT streaming response part " + string(i) + "\n"))

			time.Sleep(500 * time.Millisecond)
		}
	}()
	return reader, nil
}
