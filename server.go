package main

import (
	"fmt"
	"net"
	"net/rpc"
)

var fileList []ServFile

// Server struct
type Server struct{}

// RequestServFileList Gets the list of file infos of the files being served
func (serv *Server) RequestServFileList(input string, reply *[]ServFile) error {
	*reply = make([]ServFile, len(fileList), len(fileList))
	copy(*reply, fileList)
	return nil
}

// RequestFile Gets the file from server
func (serv *Server) RequestFile(input string, reply *[]byte) error {
	return nil
}

func server(port string) {
	fmt.Println("Gosync server starting at port", port)
	// create a list of ServFiles
	fileList = CreateServFileList(".")
	// Print the serv file list
	fmt.Println("Files being served:")
	for _, file := range fileList {
		fmt.Println(file.Name, ":", file.Sha1)
	}

	rpc.Register(new(Server))
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		c, err := listener.Accept()
		if err != nil {
			continue
		}
		remoteAddr := c.RemoteAddr()
		fmt.Println("Accepted connection from", remoteAddr)
		go rpc.ServeConn(c)
	}
}
