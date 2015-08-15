package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	//"path/filepath"
)

// ServFile stores the filename and sha1 of a file
type ServFile struct {
	Name string
	Sha1 string
}

// CreateServFileList creates a ServFile list for all the files in the dir
func CreateServFileList(dirname string) []ServFile {
	var fileList []ServFile
	// get the files in the directory
	// files, _ := filepath.Glob("*")
	files, _ := ioutil.ReadDir(dirname)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filename := file.Name()
		fileBytes, err := GetFileBytes(filename)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fileList = append(fileList, ServFile{filename, GetHash(fileBytes)})
	}
	return fileList
}

// GetFileBytes gets the bytes of a files
func GetFileBytes(filename string) ([]byte, error) {
	filebytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return []byte{}, err
	}
	return filebytes, err
}

// GetHash gets the sha1 hash of a file
func GetHash(filebytes []byte) string {
	h := sha1.New()
	h.Write(filebytes)
	return hex.EncodeToString(h.Sum([]byte{}))
}
