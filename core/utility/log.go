package utility

import (
	"fmt"
	"strings"
)

func Log(message string) {
	fmt.Println("-----> ", message)
}

func LogCouldNotExecuteCommand(reasons []string) {
	fmt.Printf("Could not execute command due to the following reason(s):\n%s", strings.Join(reasons, "\n"))
}