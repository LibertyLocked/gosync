package main

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestCompression(t *testing.T) {
	input := "input for compression test"
	compressed := Compress([]byte(input))
	decompressed, err := Decompress(compressed)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	output := string(decompressed)
	if input != output {
		t.Fail()
	}
}

func TestEncryption(t *testing.T) {
	input := "input for encryption test"
	key := GetKeyHash("mySecretKey")
	encrypted, err := Encrypt([]byte(input), key)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	decrypted, err := Decrypt(encrypted, key)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	output := string(decrypted)
	if input != output {
		t.Fail()
	}
}

func TestGetFileHash(t *testing.T) {
	sha1str := GetFileHash([]byte("testinput"))
	if sha1str != "5f2a30a0b1ea8ea0f1b0ce8c52338ed940334344" {
		t.Fail()
	}
}

func TestGetKeyHash(t *testing.T) {
	sha256bytes := GetKeyHash("testinput")
	if hex.EncodeToString(sha256bytes) != "e0b759f336aefd2ff5b31534f23d98cedfdca407850ba5c6c99502c424441ab7" {
		t.Fail()
	}
}
