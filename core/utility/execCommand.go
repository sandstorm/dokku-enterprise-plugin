package utility

import (
	"os/exec"
	"strings"
)

/**
 * Execute a command and return the output as string
 */
func ExecCommand(commandParts ...string) (string) {
	command := commandParts[0]
	commandArgs := commandParts[1:]
	cmd := exec.Command(command, commandArgs...)
	output, err := cmd.CombinedOutput()

	if (err != nil) {
		return ""
	} else {
		return strings.TrimSpace(string(output[:]))
	}
}
