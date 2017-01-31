package utility

import (
	"path/filepath"
	"os"
	"fmt"
)

// Copies all files from source to target directory
// But use with caution: Target directory is overriden completely!
func CopyAndOverrideDirectory(source, target string) {
	source = filepath.Clean(source)
	target = filepath.Clean(target)

	os.MkdirAll(filepath.Dir(target), 0777)

	// we have to make sure that target does not exist
	// -> otherwise the complete source folder will be copied into target (instead of only its contents)
	err := os.RemoveAll(target)
	if err != nil {
		fmt.Errorf("Could not remove persistent-data-folder %s. Error was: %v", target, err)
	}

	ExecCommandAndFailWithFatalErrorOnError("cp", "-R", source, target)
}