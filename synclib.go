package main

import (
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"strings"
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
		// Exclude directories
		if file.IsDir() {
			continue
		}
		// Exclude gosync itself
		if strings.HasPrefix(strings.ToLower(file.Name()), strings.ToLower(GoSyncExeName)) {
			continue
		}
		filename := file.Name()
		fileBytes, err := GetFileBytes(filename)
		if err != nil {
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

// Diff compares two ServFile slices and returns the info of files to add and to delete
func Diff(local, remote []ServFile) (toAdd, toDel []ServFile) {
	m := make(map[ServFile]int)
	for _, localFile := range local {
		m[localFile]++
	}
	for _, remoteFile := range remote {
		m[remoteFile] += 2
	}
	for mKey, mVal := range m {
		if mVal == 1 {
			// this file is only on local
			toDel = append(toDel, mKey)
		} else if mVal == 2 {
			// this file is only on remote
			toAdd = append(toAdd, mKey)
		}
	}
	return
}
