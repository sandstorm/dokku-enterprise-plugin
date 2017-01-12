package utility

import (
	"os/exec"
	"strings"
	"io"
	"os"
	"log"
	"path/filepath"
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

func ExecCommandAndFailWithFatalErrorOnError(commandParts ...string) {
	command := commandParts[0]
	commandArgs := commandParts[1:]
	cmd := exec.Command(command, commandArgs...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("FATAL: There was an error running command %v; error was: %v", commandParts, err)
	}
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


func ExecCommandAndDumpResultToFile(resultFile string, commandParts ...string) {
	command := commandParts[0]
	commandArgs := commandParts[1:]
	cmd := exec.Command(command, commandArgs...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("ERROR: Stdout Pipe could not be redirected: %v", err)
	}

	os.MkdirAll(filepath.Dir(resultFile), os.ModePerm)
	resultFileHandle, err := os.Create(resultFile)
	if err != nil {
		log.Fatalf("ERROR: Result file %s could not be created: %v", resultFile, err)
	}

	cmd.Start()

	_, err = io.Copy(resultFileHandle, stdout)
	if err != nil {
		log.Fatalf("ERROR: Result file %s could not be copied to: %v", resultFile, err)
	}
	cmd.Wait()

	stat, err := resultFileHandle.Stat()
	if err != nil {
		log.Fatalf("ERROR: Result file stats for %s could not be obtained: %v", resultFile, err)
	}

	if stat.Size() == 0 {
		log.Fatalf("ERROR: File %s is empty! Something has gone wrong!!", resultFile)
	}
	resultFileHandle.Close()
}