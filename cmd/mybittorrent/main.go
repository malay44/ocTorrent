package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"

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
		conn, err := net.Dial("tcp", peer)
		if err != nil {
			fmt.Println(err)
			return
		}
		var buf []byte
		buf = append(buf, 19)
		buf = append(buf, []byte("BitTorrent protocol")...)
		buf = append(buf, make([]byte, 8)...)
		buf = append(buf, torrent.Info.hash...)
		buf = append(buf, []byte("00112233445566778899")...)
		_, err = conn.Write(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		answer := make([]byte, 68)
		io.ReadFull(conn, answer)
		// Read last 20 bytes (peerID)
		fmt.Printf("Peer ID: %x\n", answer[48:])	
	}else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
