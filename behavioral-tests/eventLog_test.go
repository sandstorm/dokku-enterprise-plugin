package main

import (
	"github.com/sandstorm/dokku-enterprise-plugin/core/utility"
	"github.com/DATA-DOG/godog"
)

func theEventLogIsEmpty() error {
	utility.ExecCommand("ssh", "dokku@dokku.me", "rm", "-Rf", "/home/dokku/.event-log-tmp/")
	return nil
}

func iExpectEventLogEntry(amount int) error {
	utility.ExecCommand("ssh", "dokku@dokku.me", "ls", "/home/dokku/.event-log-tmp/")
	return godog.ErrPending
}
