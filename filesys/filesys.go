package filesys

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

/*ReadFile performs checks and retrieves contents of file*/
func ReadFile(path string) ([]byte, error) {
	// Checking if file exists
	_, statErr := os.Stat(path)
	if statErr != nil {
		if os.IsNotExist(statErr) {
			return nil, errors.New("file does not exist")
		}
		return nil, statErr
	}

	// test permissions
	file, openErr := os.OpenFile(path, os.O_RDWR, 0666)
	if openErr != nil {
		if os.IsPermission(openErr) {
			return nil, errors.New("write permission denied")
		}
		if os.IsPermission(openErr) {
			return nil, errors.New("read permission denied")
		}
		file.Close()
		return nil, openErr
	}

	defer file.Close()

	data, readErr := ioutil.ReadAll(file)
	if readErr != nil {
		return nil, readErr
	}

	return data, nil
}

/*WriteFile saves data to the file*/
func WriteFile(data []byte, path string) error {
	newFile, createErr := os.Create(path)
	if createErr != nil {
		log.Fatal(createErr)
	}

	defer newFile.Close()

	_, writeErr := newFile.Write(data)
	if writeErr != nil {
		log.Fatal(writeErr)
	}

	return nil
}

/*Replace updates contents of the file*/
func Replace(path string, data []byte) error {
	newpath := filepath.Dir(path) + "_" + filepath.Base(path)
	writeErr := WriteFile(data, newpath)
	if writeErr != nil {
		return writeErr
	}

	removeErr := os.Remove(path)
	if removeErr != nil {
		return removeErr
	}

	renameErr := os.Rename(newpath, path)
	if renameErr != nil {
		return renameErr
	}

	return nil
}
