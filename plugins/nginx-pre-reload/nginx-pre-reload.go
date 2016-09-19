package main

import (
	"github.com/sandstorm/dokku-enterprise-plugin/core/dokku"
	"github.com/sandstorm/dokku-enterprise-plugin/core/utility"
	"io/ioutil"
	"os"
)

// http://dokku.viewdocs.io/dokku/development/plugin-triggers/#nginx-pre-reload
func main() {
	app := os.Args[1]

	utility.Log("Running dokku-enterprise nginx-pre-reload hook")

	nginxConfDirectory := dokku.DokkuRoot() + "/" + app + "/nginx.conf.d"
	os.RemoveAll(nginxConfDirectory)

	appContainerId := dokku.GetAppContainerId(app)
	utility.ExecCommand("docker", "cp", appContainerId+":/app/nginx.conf.d", nginxConfDirectory)

	files, _ := ioutil.ReadDir(nginxConfDirectory)
	for _, f := range files {
		utility.Log("custom nginx.conf.d: including " + f.Name())
	}

}
