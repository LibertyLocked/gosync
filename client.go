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
	// Get the file infos of the files being served
	err = c.Call("Server.RequestServFileList", "", &remoteFileList)
	if err != nil {
		fmt.Println(err)
	} else {
		// successfully obtained the list of file infos
		// diff once to find the files we need to download
		toAdd, _ := Diff(localFileList, remoteFileList)
		for _, val := range toAdd {
			// download the files needed
			fmt.Println("+", val.Name, ":", val.Sha1)
			var buffer []byte
			c.Call("Server.RequestFile", val, &buffer)
			ioutil.WriteFile(val.Name, buffer, os.FileMode(0644))
		}
		// refresh local file list
		localFileList = CreateServFileList(".")
		// diff again to find the files we need to remove
		_, toDel := Diff(localFileList, remoteFileList)
		for _, val := range toDel {
			// delete the files that aren't needed
			fmt.Println("-", val.Name, ":", val.Sha1)
			os.Remove(val.Name)
		}
	}
	c.Close()
}
