package main

import (
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
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
		fileList = append(fileList, ServFile{filename, GetFileHash(fileBytes)})
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

// GetFileHash gets the sha1 hash of a file
func GetFileHash(filebytes []byte) string {
	h := sha1.New()
	h.Write(filebytes)
	return hex.EncodeToString(h.Sum([]byte{}))
}

// GetKeyHash gets the sha256 hash of a string
func GetKeyHash(str string) []byte {
	h := sha256.New()
	h.Write([]byte(str))
	return h.Sum([]byte{})
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

// Encrypt encrypts plaintext to ciphertext using AES
func Encrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(plaintext)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

// Decrypt decrypts ciphertext to plaintext using AES
func Decrypt(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(ciphertext, ciphertext)
	data, err := base64.StdEncoding.DecodeString(string(ciphertext))
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Compress uses zlib to compress a byte array
func Compress(data []byte) []byte {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(data)
	w.Flush()
	w.Close()
	return b.Bytes()
}

// Decompress uses zlib to decompress a compressed byte array
func Decompress(compressed []byte) ([]byte, error) {
	var out bytes.Buffer
	in := bytes.NewBuffer(compressed)
	r, err := zlib.NewReader(in)
	io.Copy(&out, r)
	r.Close()
	return out.Bytes(), err
}
