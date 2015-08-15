package main

import (
	"fmt"
	"net/rpc"
)

func client(serverAddr string) {
	fmt.Println("Connecting to", serverAddr)
	c, err := rpc.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	var fileList []ServFile
	// Get the file infos of the files being served
	err = c.Call("Server.RequestServFileList", "", &fileList)
	if err != nil {
		fmt.Println(err)
	} else {
		// successfully obtained the list of file infos
		fmt.Println("Remote hashes:")
		for _, value := range fileList {
			fmt.Println(value.Name, ":", value.Sha1)
		}
	}
	c.Close()
}
