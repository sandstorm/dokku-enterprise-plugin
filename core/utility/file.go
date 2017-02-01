package utility

import (
	"path/filepath"
	"os"
	"fmt"
	"io"
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

func FileExists(filepath string) bool {
	_, err := os.Stat(filepath)

	if err == nil {
		return true
	}

	return ! os.IsNotExist(err)
}

func DirectoryIsEmpty(filepath string) bool {
	f, err := os.Open(filepath)
	if err != nil {
		return false
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true
	}
	return false
}