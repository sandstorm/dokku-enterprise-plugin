package dokku

import (
	"os"
)


func DokkuRoot() (string) {
	return os.Getenv("DOKKU_ROOT");
}