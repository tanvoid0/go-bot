package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Read(message string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(message)
	query, _ := reader.ReadString('\n')
	return strings.TrimSpace(query)
}

func ClearScreen() {
	//cmd := exec.Command("clear")
	//err := cmd.Run()
	//if err != nil {
	//	fmt.Println("Failed to clear screen:", err)
	//}
	fmt.Printf("\033[1A\033[K")

	//fmt.Print("\033[2J\033[H")
}
