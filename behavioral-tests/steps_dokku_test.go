package main

import (
	"fmt"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/sandstorm/dokku-enterprise-plugin/core/utility"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"github.com/sandstorm/dokku-enterprise-plugin/behavioral-tests/jsonQueryHelper"
)

// Create a quite-minimal dokku Dockerfile application (as basis for testing)
// which just delivers a static file.
func iHaveAnEmptyDockerfileApplication() error {
	os.RemoveAll("/tmp/bdd-test-app")
	os.Mkdir("/tmp/bdd-test-app", 0755)
	os.Chdir("/tmp/bdd-test-app")

	dockerfileContents := `
		FROM nginx:stable-alpine
		ADD . /app
		EXPOSE 5000
		RUN cp /app/nginx.conf /etc/nginx/nginx.conf
		CMD ["nginx", "-g", "daemon off;"]
	`
	ioutil.WriteFile("Dockerfile", []byte(dockerfileContents), 0644)

	checksContents := `
		WAIT=1
		http://test/ welcome
	`
	ioutil.WriteFile("CHECKS", []byte(checksContents), 0644)

	nginxConf := `
		events {
			worker_connections  1024;
		}
		http {
			server {
				listen 5000;
				listen 80;
				listen 443;
				location / {
					root   /app;
					index  index.html index.htm;
				}
			}
		}
	`
	ioutil.WriteFile("nginx.conf", []byte(nginxConf), 0644)

	index := "welcome"
	ioutil.WriteFile("index.html", []byte(index), 0644)

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
	exec.Command("git", "init").Run()
	exec.Command("git", "add", ".").Run()
	exec.Command("git", "commit", "-m", "Initial commit").Run()
	result := utility.ExecCommand("git", "push", "--force", "dokku@dokku.me:"+applicationName, "HEAD:master")

	if strings.Contains(result, "Application deployed:") {
		return nil
	} else {
		return fmt.Errorf("Deployment was not successful. Full log: %s", result)
	}
}

var dokkuResponseBody string

// Call dokku with some arguments
func iCallDokku(dokkuArguments string) error {
	args := strings.Split(dokkuArguments, " ")

	dokkuResponseBody = utility.ExecCommand(append([]string{"ssh", "dokku@dokku.me"}, args...)...)

	return nil
}

func iGetBackAJSONObjectWithTheFollowingStructure(comparators *gherkin.DataTable) error {
	return jsonQueryHelper.AssertJsonStructure(dokkuResponseBody, comparators)
}

// the HTTP response body as string; filled as result of iCallTheURLOfTheApplication()
var httpResponseBodyAfterCallingApplicationUrl string

// Call the URL of an application.
func iCallTheURLOfTheApplication(urlPath, applicationName string) error {

	// First, we need to figure out the port the application is running on; and we need to prefix "dokku.me"
	// for the host to work properly across Vagrant development environments.

	domainNames := utility.ExecCommand("ssh", "dokku@dokku.me", "urls", applicationName)
	domainNamesArray := strings.Split(domainNames, "\n")
	parsedUrl, _ := url.Parse(domainNamesArray[0])
	splitHost := strings.Split(parsedUrl.Host, ":")
	parsedUrl.Host = "dokku.me:" + splitHost[1]
	parsedUrl.Path = urlPath

	// Do the actual request, and parse the response
	result, err := http.Get(parsedUrl.String())

	if err != nil {
		return fmt.Errorf("Error while calling %s: %s", parsedUrl.String(), err)
	}

	defer result.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(result.Body)
	httpResponseBodyAfterCallingApplicationUrl = string(bodyBytes[:])

	return nil
}

func theResponseShouldContain(substring string) error {
	if !strings.Contains(httpResponseBodyAfterCallingApplicationUrl, substring) {
		return fmt.Errorf("String '%s' should contain '%s', but did not.", httpResponseBodyAfterCallingApplicationUrl, substring)
	}
	return nil
}

func theResponseShouldNotContain(substring string) error {
	if strings.Contains(httpResponseBodyAfterCallingApplicationUrl, substring) {
		return fmt.Errorf("String '%s' should not contain '%s', but did.", httpResponseBodyAfterCallingApplicationUrl, substring)
	}
	return nil
}
