package common

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

// StringInSlice - checks if specified string is in array
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

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

// CompressBytes - compress given byte array via gzip
func CompressBytes(b []byte) []byte {
	var bBuff bytes.Buffer
	gz, err := gzip.NewWriterLevel(&bBuff, gzip.BestCompression)

	if err != nil {
		panic(err)
	}

	if _, err := gz.Write(b); err != nil {
		panic(err)
	}

	if err := gz.Flush(); err != nil {
		panic(err)
	}

	if err := gz.Close(); err != nil {
		panic(err)
	}

	return bBuff.Bytes()
}

// DecompressBytes - decompress given bytes via gzip
func DecompressBytes(b []byte) ([]byte, error) {
	var bBuff bytes.Buffer
	bBuff.Write(b)
	r, err := gzip.NewReader(&bBuff)

	if err != nil {
		return nil, err
	}

	r.Close()

	s, err := ioutil.ReadAll(r)
	return s, nil
}

// GetCurrentDir - returns current execution directory
func GetCurrentDir() string {
	currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return currentDir
}

// GetCurrentTime - Fetch current UTC time
func GetCurrentTime() time.Time {
	return time.Now().UTC()
}
