package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/codecrafters-io/bittorrent-starter-go/bencode"
)

func main() {
	command := os.Args[1]
	if command == "decode" {
		bencodedValue := os.Args[2]
		
		decoded,_, err := bencode.Decode(bencodedValue, 0)
		if err != nil {
			fmt.Println(err)
			return
		}
		
		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))
	} else if command == "info" {
		torrentFile := os.Args[2]
		_, err := parseTorrentFile(torrentFile)
		if err != nil {
			fmt.Println(err)
			return
		}
	}else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
