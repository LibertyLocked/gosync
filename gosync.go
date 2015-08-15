package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 2 && os.Args[1] == "-s" {
		// run it in server mode
		port := os.Args[2]
		server(port)
	} else if len(os.Args) > 2 && os.Args[1] == "-c" {
		// server address is in args
		serverAddr := os.Args[2]
		client(serverAddr)
	} else if len(os.Args) == 1 {
		// no args. run in client mode and ask for user input
		var serverAddr string
		fmt.Println("Please enter server address (e.g. localhost:9999): ")
		fmt.Scanln(&serverAddr)
		client(serverAddr)
	} else {
		fmt.Println("Invalid arguments!")
		fmt.Println("For servers: gosync -s <port>")
		fmt.Println("For clients: gosync -c <addr:port>")
	}
}
