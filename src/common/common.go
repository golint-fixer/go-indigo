package common

import (
	"encoding/gob"
	"log"
	"os"
	"path/filepath"
)

// WriteGob - create gob from specified object, at filePath
func WriteGob(filePath string, object interface{}) error {
	file, err := os.Create(filePath)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

// ReadGob - read gob specified at path
func ReadGob(filePath string, object interface{}) error {
	file, err := os.Open(filePath)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}

// GetCurrentDir - returns current execution directory
func GetCurrentDir() string {
	currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return currentDir
}
