package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/codecrafters-io/bittorrent-starter-go/bencode"
)

type Torrent struct {
	Announce string
	
	Info     struct {
		Length      int
		PieceLength int
		Pieces      []string
		hash 		[]byte
	}
}
func (torrent *Torrent) Print () {
	fmt.Println("Tracker URL:", torrent.Announce)
	fmt.Println("Length:", torrent.Info.Length)
	fmt.Println("Info Hash: " + hex.EncodeToString(torrent.Info.hash))
	fmt.Println("Piece Length:", torrent.Info.PieceLength)
	fmt.Println("Piece Hashes:")
	for _, piece := range torrent.Info.Pieces {
		fmt.Println(hex.EncodeToString([]byte(piece)))
	}
}

func parseTorrentFile(torrentFileName string, torrent *Torrent) (map[string]interface{}, error) {
	file, err := os.ReadFile(torrentFileName)
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
	torrent.Announce = torrentData["announce"].(string)


	info, ok := torrentData["info"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid torrent file, info key not found")
	}
	torrent.Info.Length = info["length"].(int)

	text, err := bencode.Encode(info)
	if err != nil {
		return nil, err
	}
	var sha = sha1.New()
	sha.Write([]byte(text))
	var encrypted = sha.Sum(nil)
	torrent.Info.hash = encrypted

	torrent.Info.PieceLength = info["piece length"].(int)

	// Piece Hashes
	pieceHashes := info["pieces"].(string)
	pieceHashesLength := len(pieceHashes)
	const pieceLength = 20
	pieces := pieceHashesLength / pieceLength

	torrent.Info.Pieces = make([]string, pieces)
	for i := 0; i < pieces; i++ {
		start := i * pieceLength
		end := (i + 1) * pieceLength
		torrent.Info.Pieces[i] = pieceHashes[start:end]
	}


	return torrentData, nil
}