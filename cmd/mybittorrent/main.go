package main

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/codecrafters-io/bittorrent-starter-go/bencode"
)

const Debug bool = false

func getPeers(trackerURL string, infoHash []byte, peerID string, left int) ([]string, error) {
	trackerURLParsed, err := url.Parse(trackerURL)
	if err != nil {
		return nil, err
	}

	query := trackerURLParsed.Query()
	query.Set("info_hash", string(infoHash))
	query.Set("peer_id", peerID)
	query.Set("port", "6881")
	query.Set("uploaded", "0")
	query.Set("downloaded", "0")
	query.Set("left", strconv.Itoa(left))
	query.Set("compact", "1")
	trackerURLParsed.RawQuery = query.Encode()

	if Debug {
		fmt.Println("Tracker URL:")
		fmt.Println(trackerURLParsed.String())
	}

	httpResponse, err := http.Get(trackerURLParsed.String())
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}

	// Decode the response body
	decodedResponse, _, err := bencode.Decode(string(responseBody), 0)
	if err != nil {
		return nil, err
	}

	// parse the response body to json if debug is enabled
	if Debug {
		jsonOutput, _ := json.MarshalIndent(decodedResponse, "", "  ")
		fmt.Println("\nTracker Response:")
		fmt.Println(string(jsonOutput))
	}

	// Get the peers from the response
	peers, ok := decodedResponse.(map[string]interface{})["peers"].(string)
	if !ok {
		return nil, fmt.Errorf("peers key not found in the response")
	}

	// Parse the peers
	peersList := make([]string, 0)
	for i := 0; i < len(peers); i += 6 {
		port := binary.BigEndian.Uint16([]byte{peers[i+4], peers[i+5]})
		ip := fmt.Sprintf("%d.%d.%d.%d", peers[i], peers[i+1], peers[i+2], peers[i+3])
		peersList = append(peersList, fmt.Sprintf("%s:%d", ip, port))
	}

	return peersList, nil
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
		peerList, err := getPeers(torrent.Announce, torrent.Info.hash, "lWM8BIeMZhfdHjGgLHBS", torrent.Info.Length)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, peer := range peerList {
			fmt.Println(peer)
		}
	}else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
