package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GoSyncExeName the name of the gosync executable
var GoSyncExeName string

// FlagRm flag indicates wheter we should delete files that aren't sync'd
var FlagRm bool

func main() {
	GoSyncExeName = filepath.Base(os.Args[0])

	// flags
	serverMode := false
	FlagRm = false
	addr := ""

	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-") {
			// args that start with '-'
			switch arg {
			case "-c":
				serverMode = false
			case "-s":
				serverMode = true
			case "-rm":
				FlagRm = true
			case "-help":
				printHelp()
				return
			default:
				fmt.Println("Unknown option:", arg)
				return
			}
		} else {
			// args don't start with '-'
			addr = arg
		}
	}

	if addr == "" {
		if serverMode {
			fmt.Println("Please enter port number (e.g. 9999): ")
		} else {
			fmt.Println("Please enter server address (e.g. localhost:9999): ")
		}
		fmt.Scanln(&addr)
	}

	if serverMode {
		server(addr)
	} else {
		client(addr)
	}
}

func printHelp() {
	fmt.Println("Usage:", GoSyncExeName, "[-Options] [port/addr:port]")
	fmt.Println("\nFor servers:", GoSyncExeName, "-s [port]")
	fmt.Println("\tExample:", GoSyncExeName, "-s 9999")
	fmt.Println("For clients:", GoSyncExeName, "-c [addr:port]")
	fmt.Println("\tExample:", GoSyncExeName, "-c localhost:9999")
	fmt.Println("Other flags:")
	fmt.Println("-rm\tRemove out-of-sync local files")
}
