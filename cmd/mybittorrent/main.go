package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/codecrafters-io/bittorrent-starter-go/bencode"
)

func getPeers(trackerURL string, infoHash []byte, peerID string) {
	fmt.Println("Getting peers from tracker")
	fmt.Println("Tracker URL: " + trackerURL)
	fmt.Println("Info Hash: " + hex.EncodeToString(infoHash))
	fmt.Println("Peer ID: " + peerID)
}

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
		var torrent Torrent
		_, err := parseTorrentFile(torrentFile, &torrent)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Tracker URL:", torrent.Announce)
		fmt.Println("Length:", torrent.Info.Length)
		fmt.Println("Info Hash: " + hex.EncodeToString(torrent.Info.hash))
		fmt.Println("Piece Length:", torrent.Info.PieceLength)
		fmt.Println("Piece Hashes:")
		for _, piece := range torrent.Info.Pieces {
			fmt.Println(hex.EncodeToString([]byte(piece)))
		}
	} else if command == "peers" {
		torrentFile := os.Args[2]
		var torrent Torrent
		_, err := parseTorrentFile(torrentFile, &torrent)
		if err != nil {
			fmt.Println(err)
			return
		}
		getPeers(torrent.Announce, torrent.Info.hash, "12345678901234567890")
	}else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
