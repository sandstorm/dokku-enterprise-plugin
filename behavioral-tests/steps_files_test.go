package main

import (
	"github.com/sandstorm/dokku-enterprise-plugin/core/utility"
	"fmt"
	"regexp"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/DATA-DOG/godog"
)

func anEmptyFolderExists(folder string) error {
	utility.ExecCommand("ssh", "root@dokku.me", "rm", "-Rf", folder)
	utility.ExecCommand("ssh", "root@dokku.me", "mkdir", "-p", folder)
	utility.ExecCommand("ssh", "root@dokku.me", "chmod", "-R", "777", folder)
	return nil
}

func iExpectAFileInFolder(filePattern, folder string) error {
	result := utility.ExecCommand("ssh", "root@dokku.me", "ls", folder)

	matched, err := regexp.MatchString(filePattern, result)
	if (err != nil) {
		return fmt.Errorf("ERROR while regex: %v", err)
	}

	if (!matched) {
		return fmt.Errorf("pattern not found. Files: %v", result)
	}
	return nil
}


func iExpectAFileWithContents(arg1 string, arg2 *gherkin.DocString) error {
	return godog.ErrPending
}

func aFileIsCreatedWithContents(arg1 string, arg2 *gherkin.DocString) error {
	return godog.ErrPending
}