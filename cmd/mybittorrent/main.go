package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/codecrafters-io/bittorrent-starter-go/bencode"
)

const Debug bool = false

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
		peerList, err := retrievePeers(torrent.Announce, torrent.Info.hash, "lWM8BIeMZhfdHjGgLHBS", torrent.Info.Length)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, peer := range peerList {
			fmt.Println(peer)
		}
	} else if command == "handshake" {
		torrentFile := os.Args[2]
		peer := os.Args[3]
		var torrent Torrent
		_, err := parseTorrentFile(torrentFile, &torrent)
		if err != nil {
			fmt.Println(err)
			return
		}
		_,answer, err := Handshake(peer, &torrent.Info.hash)
		if err != nil {
			fmt.Println(err)
			return
		}
		// Read last 20 bytes (peerID)
		fmt.Printf("Peer ID: %x\n", answer[48:])	
	} else if command == "download_piece" {
		// $ ./your_bittorrent.sh download_piece -o /tmp/test-piece-0 sample.torrent 0
		outputFile := os.Args[3]
		torrentFile := os.Args[4]
		pieceIndex := os.Args[5]
		index, err := strconv.Atoi(pieceIndex)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = RetrievePiece(outputFile, torrentFile, index)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
