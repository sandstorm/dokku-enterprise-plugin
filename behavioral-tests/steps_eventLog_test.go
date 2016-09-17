package main

import (
	"github.com/sandstorm/dokku-enterprise-plugin/core/utility"
	"strconv"
	"fmt"
)

func theEventLogIsEmpty() error {
	utility.ExecCommand("ssh", "root@dokku.me", "rm", "-Rf", "/home/dokku/.event-log-tmp/")
	return nil
}

func iExpectEventLogEntry(expectedNumberOfLines int) error {
	result := utility.ExecCommand("ssh", "root@dokku.me", "ls /home/dokku/.event-log-tmp/ | wc -l")
	numberOfLines, _ := strconv.Atoi(result)

	if (numberOfLines != expectedNumberOfLines) {
		return fmt.Errorf("Expected %d number of log entries, got %d", expectedNumberOfLines, numberOfLines)
	}

	return nil
}
