package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/codecrafters-io/bittorrent-starter-go/bencode"
)

func retrievePeers(trackerURL string, infoHash []byte, peerID string, left int) ([]string, error) {
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

	peersList := make([]string, 0)
	// Parse the peers
	for i := 0; i < len(peers); i += 6 {
		port := binary.BigEndian.Uint16([]byte{peers[i+4], peers[i+5]})
		ip := fmt.Sprintf("%d.%d.%d.%d", peers[i], peers[i+1], peers[i+2], peers[i+3])
		peersList = append(peersList, fmt.Sprintf("%s:%d", ip, port))
	}

	return peersList, nil
}