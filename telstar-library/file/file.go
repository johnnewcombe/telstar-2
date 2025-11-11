package file

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/johnnewcombe/telstar-library/logger"
)

// Exists checks for the existance of a file at the specified path.
func Exists(path string) bool {

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
func DeleteFiles(fileSpec string) {

	var (
		fileMatches []string
		err         error
	)
	// delete the output directory of json files
	if fileMatches, err = filepath.Glob(fileSpec); err != nil {
		logger.LogError.Print(err)
		return
	}
	for _, f := range fileMatches {
		if err = os.Remove(f); err != nil {
			logger.LogError.Print(err)
			return
		}
	}
}

func WriteFile(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0644)
}

func ReadFile(path string) ([]byte, error) {

	var (
		bytes []byte
		err   error
	)
	if bytes, err = os.ReadFile(path); err != nil {
		return bytes, err
	}
	return bytes, nil
}
