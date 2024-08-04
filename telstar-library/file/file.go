package file

import (
	"errors"
	"github.com/johnnewcombe/telstar-library/logger"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
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
	return ioutil.WriteFile(filename, data, 0644)
}

func ReadFile(path string) ([]byte, error) {

	var (
		bytes []byte
		err   error
	)
	if bytes, err = ioutil.ReadFile(path); err != nil {
		return bytes, err
	}
	return bytes, nil
}

func ReadFiles(directory string) ([]fs.FileInfo, error) {

	var (
		err   error
		files []fs.FileInfo
	)

	files, err = ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	return files, nil
}
