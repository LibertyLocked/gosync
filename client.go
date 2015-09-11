package main

import (
	"fmt"
	"io/ioutil"
	"net/rpc"
	"os"
)

var localFileList []ServFile

func client(serverAddr string) {
	fmt.Println("Connecting to", serverAddr)
	// Create list of ServFiles in local directory
	localFileList = CreateServFileList(".")
	c, err := rpc.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	var remoteFileList []ServFile
	// Ask server if encryption is used
	var encryptionMode bool
	err = c.Call("Server.RequestEncryptionMode", "", &encryptionMode)
	if !encryptionMode {
		UseEncryption = false
	} else if encryptionMode && !UseEncryption {
		fmt.Println("Sync aborted: Server requires encryption")
		return
	}
	// Ask server if compression is used
	err = c.Call("Server.RequestCompressionMode", "", &UseCompression)
	// Get the file infos of the files being served
	err = c.Call("Server.RequestServFileList", "", &remoteFileList)
	if err != nil {
		fmt.Println(err)
	} else {
		// successfully obtained the list of file infos
		// diff once to find the files we need to download
		toAdd, _ := Diff(localFileList, remoteFileList)
		fmt.Println("Downloading", len(toAdd), "file(s)")
		for _, file := range toAdd {
			// download the files needed
			var buffer []byte
			c.Call("Server.RequestFile", file, &buffer)
			// Decrypt if encryption is used
			if UseEncryption {
				buffer, err = Decrypt(buffer, AESKey)
				if err != nil {
					fmt.Println("Error decrypting file:", file.Name)
					continue
				}
			}
			// Decompress if compression is used
			if UseCompression {
				buffer, err = Decompress(buffer)
				if err != nil {
					fmt.Println("Error decompressing file:", file.Name)
					continue
				}
			}
			ioutil.WriteFile(file.Name, buffer, os.FileMode(0644))
			// tell user the file has been downloaded
			fmt.Println("+", file.Name, ":", file.Sha1)
			// check sha1 of the written file to ensure integrity
			writtenBytes, _ := GetFileBytes(file.Name)
			if GetFileHash(writtenBytes) != file.Sha1 {
				fmt.Println("Hash mismatch! Downloaded file is corrupted!")
			}
		}

		if FlagRm {
			// refresh local file list
			localFileList = CreateServFileList(".")
			// diff again to find the files we need to remove (if FlagRm is true)
			_, toDel := Diff(localFileList, remoteFileList)
			fmt.Println("Removing", len(toDel), "file(s)")
			for _, val := range toDel {
				// delete the files that aren't needed
				os.Remove(val.Name)
				// tell user the file has been deleted
				fmt.Println("-", val.Name, ":", val.Sha1)
			}
		}
	}
	fmt.Println("Sync completed! Closing connection")
	c.Close()
}
