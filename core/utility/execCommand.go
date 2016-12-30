package utility

import (
	"os/exec"
	"strings"
	"io"
)

/**
 * Execute a command and return the output as string
 */
func ExecCommand(commandParts ...string) string {
	return execCommandInternal(nil, commandParts...)
}

func ExecCommandWithStdIn(stdin io.Reader, commandParts ...string) string {
	return execCommandInternal(stdin, commandParts...)
}

/**
 * Execute a command and return the output as string
 */
func execCommandInternal(stdin io.Reader, commandParts ...string) string {
	command := commandParts[0]
	commandArgs := commandParts[1:]
	cmd := exec.Command(command, commandArgs...)
	if stdin != nil {
		cmd.Stdin = stdin
	}
	output, err := cmd.CombinedOutput()

	if err != nil {
		return "ERROR running command: " + strings.Join(commandParts, " ") + " - " + strings.TrimSpace(string(output[:]))
	} else {
		return strings.TrimSpace(string(output[:]))
	}
}
