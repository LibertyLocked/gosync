package main

import (
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
