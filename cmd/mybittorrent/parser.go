package main

import (
	"crypto/sha1"
	"fmt"
	"os"

	"github.com/codecrafters-io/bittorrent-starter-go/bencode"
)

func parseTorrentFile(torrentFile string) (map[string]interface{}, error) {
	file, err := os.ReadFile(torrentFile)
	if err != nil {
		return nil, err
	}

	bencodedSting := string(file)
	data, _, err := bencode.Decode(bencodedSting, 0)
	if err != nil {
		return nil, err
	}

	// if the data is a dictionary, then it is a valid torrent file
	var torrentData map[string]interface{}
	var ok bool
	if torrentData, ok = data.(map[string]interface{}); !ok {
		return nil, fmt.Errorf("invalid torrent file, dictionary not found")
	}
	fmt.Println("Tracker URL:", torrentData["announce"])

	info, ok := torrentData["info"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid torrent file, info key not found")
	}
	fmt.Println("Length:", info["length"])

	text, err := bencode.Encode(info)
	if err != nil {
		return nil, err
	}
	var sha = sha1.New()
	sha.Write([]byte(text))
	var encrypted = sha.Sum(nil)
	fmt.Println("Info Hash:", fmt.Sprintf("%x", encrypted))

	fmt.Println("Piece Length:", info["piece length"])

	// Piece Hashes
	pieceHashes := info["pieces"].(string)
	pieceHashesLength := len(pieceHashes)
	const pieceLength = 20
	pieces := pieceHashesLength / pieceLength

	fmt.Println("Piece Hashes:")
	for i := 0; i < pieces; i++ {
		start := i * pieceLength
		end := (i + 1) * pieceLength
		fmt.Printf("%x\n", pieceHashes[start:end])
	}

	// print map in json format with indent
	// jsonOutput, _ := json.MarshalIndent(torrentData, "", "  ")
	// fmt.Println(string(jsonOutput))

	return torrentData, nil
}