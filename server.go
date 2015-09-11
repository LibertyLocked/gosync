package main

import (
	"fmt"
	"net"
	"net/rpc"
)

var fileList []ServFile

// Server RPC struct
type Server struct{}

// RequestCompressionMode gets the mode of compression server is using
func (serv *Server) RequestCompressionMode(input string, reply *bool) error {
	*reply = UseCompression
	return nil
}

// RequestEncryptionMode gets the mode of encryption server is using
func (serv *Server) RequestEncryptionMode(input string, reply *bool) error {
	*reply = UseEncryption
	return nil
}

// RequestServFileList gets the list of file infos of the files being served
func (serv *Server) RequestServFileList(input string, reply *[]ServFile) error {
	fmt.Println("Client requested filelist")
	*reply = make([]ServFile, len(fileList))
	copy(*reply, fileList)
	return nil
}

// RequestFile gets the file from server
func (serv *Server) RequestFile(input ServFile, reply *[]byte) error {
	fmt.Println("Client requested file:", input.Name)
	filebytes, err := GetFileBytes(input.Name)
	if err != nil {
		fmt.Println("Error sending file:", input.Name)
		return err
	}
	// Compress if we use compression
	if UseCompression {
		filebytes = Compress(filebytes)
	}
	// Encrypt if we use encryption
	if UseEncryption {
		filebytes, err = Encrypt(filebytes, AESKey)
		if err != nil {
			fmt.Println("Error encrypting file:", input.Name)
			fmt.Println(err.Error())
			return err
		}
	}
	*reply = make([]byte, len(filebytes))
	copy(*reply, filebytes)
	return nil
}

func server(port string) {
	fmt.Println("Gosync server starting at port", port)
	// create a list of ServFiles on server
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
	fmt.Println("Awaiting connection")
	for {
		c, err := listener.Accept()
		if err != nil {
			continue
		}
		fmt.Println("Accepted connection from", c.RemoteAddr())
		go rpc.ServeConn(c)
	}
}
