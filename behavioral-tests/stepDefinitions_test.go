package main

import (
	"github.com/DATA-DOG/godog/gherkin"
	"os"
	"io/ioutil"
	"path/filepath"
	"os/exec"
	"net/url"
	"github.com/sandstorm/dokku-enterprise-plugin/core/utility"
	"net/http"
	"strings"
	"fmt"
)

func iHaveAnEmptyNodejsApplication() error {
	os.RemoveAll("/tmp/bdd-test-app")
	os.Mkdir("/tmp/bdd-test-app", 0755)
	os.Chdir("/tmp/bdd-test-app")

	packageJsonContents := `{
		"dependencies": {
			"nano-server": "*"
		},
		"scripts": {
			"start": "nano-server"
		}
	}`;
	ioutil.WriteFile("package.json", []byte(packageJsonContents), 0644)
	checksContents := "/ Listing";
	ioutil.WriteFile("CHECKS", []byte(checksContents), 0644)
	return nil
}

func iCreateTheFileWithTheFollowingContents(filePath string, contents *gherkin.DocString) error {
	os.MkdirAll(filepath.Dir(filePath), 0755)
	ioutil.WriteFile(filePath, []byte(contents.Content), 0644)

	return nil
}


func iRemoveTheFile(filePath string) error {
	os.Remove(filePath)
	return nil
}

func iDeployTheApplicationAs(applicationName string) error {
	exec.Command("git", "init").Run();
	exec.Command("git", "add", ".").Run();
	exec.Command("git", "commit", "-m", "Initial commit").Run();
	result := utility.ExecCommand("git", "push", "--force", "dokku@dokku.me:" + applicationName, "HEAD:master");

	if (strings.Contains(result, "Application deployed:")) {
		return nil
	} else {
		return fmt.Errorf("Deployment was not successful. Full log: %s", result)
	}
}

var body string;

func iCallTheURLOfTheApplication(urlPath, applicationName string) error {
	domainName := utility.ExecCommand("ssh", "dokku@dokku.me", "urls", applicationName);
	parsedUrl, _ := url.Parse(domainName)
	splitHost := strings.Split(parsedUrl.Host, ":")
	parsedUrl.Host = "dokku.me:" + splitHost[1]
	parsedUrl.Path = urlPath

	result, err := http.Get(parsedUrl.String())

	if (err != nil) {
		return fmt.Errorf("Error while calling %s: %s", parsedUrl.String(), err);
	}

	defer result.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(result.Body)
	body = string(bodyBytes[:])

	return nil
}

func theResponseShouldContain(substring string) error {
	if (!strings.Contains(body, substring)) {
		return fmt.Errorf("String '%s' should contain '%s', but did not.", body, substring)
	}
	return nil;
}

func theResponseShouldNotContain(substring string) error {
	if (strings.Contains(body, substring)) {
		return fmt.Errorf("String '%s' should not contain '%s', but did.", body, substring)
	}
	return nil;
}
