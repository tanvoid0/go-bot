package util

import (
	"github.com/charmbracelet/glamour"
	"log"
)

func MarkdownFormat(text string) string {
	rendered, err := glamour.Render(text, "dark")
	if err != nil {
		log.Fatal(err)
	}
	return rendered
}
