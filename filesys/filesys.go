package filesys

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
)

func ReadFile(path string) ([]byte, error) {
	// Checking if file exists
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("file does not exist")
		}
		return nil, err
	}
	// log.Println("File does exist. File information:")
	// log.Println(fileInfo)

	// test permissions
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		if os.IsPermission(err) {
			return nil, errors.New("write permission denied")
		}
		if os.IsPermission(err) {
			return nil, errors.New("read permission denied")
		}
		file.Close()
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func WriteFile(data []byte, path string) error {
	newFile, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}

	_, err1 := newFile.Write(data)
	if err1 != nil {
		log.Fatal(err)
	}

	newFile.Close()

	return nil
}
