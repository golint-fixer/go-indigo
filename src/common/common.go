package common

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/base64"
	"encoding/gob"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

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

// RandStringRunes - generates random string with size
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// SHA256 - hash specified byte array
func SHA256(b []byte) string {
	hash := sha256.Sum256(b)
	return base64.StdEncoding.EncodeToString(hash[:])
}

// CheckKeys - check that specified private key combinations match public key
func CheckKeys(priv string, seeds []string, pub string) bool {
	combined := seeds[0] + seeds[1]

	if SHA256([]byte(priv+combined)) == pub {
		return true
	}
	return false
}
