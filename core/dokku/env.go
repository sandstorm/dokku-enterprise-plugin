package dokku

import (
	"os"
	"io/ioutil"
)


func DokkuRoot() (string) {
	return os.Getenv("DOKKU_ROOT");
}

func Hostname() (string) {
	buffer, _ := ioutil.ReadFile("/home/dokku/HOSTNAME")
	return string(buffer)
}