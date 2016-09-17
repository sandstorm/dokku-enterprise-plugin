package dokku

import (
	"os"
	"io/ioutil"
	"strconv"
)


func DokkuRoot() (string) {
	return os.Getenv("DOKKU_ROOT");
}

func DokkuNotImplementedExit() (int) {
	exitCode, _ := strconv.Atoi(os.Getenv("DOKKU_NOT_IMPLEMENTED_EXIT"));
	return exitCode
}

func Hostname() (string) {
	buffer, _ := ioutil.ReadFile("/home/dokku/HOSTNAME")
	return string(buffer)
}