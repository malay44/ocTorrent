package main

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
)

func RetrievePiece(outputFile string, torrentFile string, pieceIndex int) error {
	fmt.Printf("Starting retrieval for piece %d from torrent file %s\n", pieceIndex, torrentFile)
	
	var torrent Torrent
	_, err := parseTorrentFile(torrentFile, &torrent)
	if err != nil {
		fmt.Printf("Error parsing torrent file: %v\n", err)
		return err
	}
	fmt.Println("Parsed torrent file successfully")
	torrent.Print()

	peerList, err := retrievePeers(torrent.Announce, torrent.Info.hash, "00112233445566778899", torrent.Info.Length)
	if err != nil {
		fmt.Printf("Error retrieving peers: %v\n", err)
		return err
	}
	fmt.Printf("Retrieved peer list: %+v\n", peerList)

	conn, _, err := Handshake(peerList[0], &torrent.Info.hash)
	if err != nil {
		fmt.Printf("Handshake failed: %v\n", err)
		return err
	}
	defer func() {
		conn.Close()
		fmt.Println("connection closed");

	}()
	fmt.Println("Handshake with peer successful")

	// Wait for bitfield message
	_, err = WaitForPeerMessage(conn, 5)
	if err != nil {
		fmt.Printf("Error waiting for bitfield message: %v\n", err)
		return err
	}
	fmt.Println("Received bitfield message")

	// Send an interested message
	fmt.Println("Sending interested message to peer")
	conn.Write([]byte{0, 0, 0, 1, 2})

	// Wait for unchoke message
	_, err = WaitForPeerMessage(conn, 1)
	if err != nil {
		fmt.Printf("Error waiting for unchoke message: %v\n", err)
		return err
	}
	fmt.Println("Received unchoke message")

	// Break the piece into blocks and send request messages
	var BLOCK_SIZE uint32 = 16 * 1024
	var byteOffset uint32 = 0
	var PieceLength uint32 = uint32(torrent.Info.PieceLength)
	if pieceIndex == len(torrent.Info.Pieces)-1 {
		PieceLength = uint32(torrent.Info.Length % torrent.Info.PieceLength)
	}
	var data = make([]byte, PieceLength)
	var count uint32 = 0

	perfectIterations := PieceLength / BLOCK_SIZE
	lastIterationLength := PieceLength % BLOCK_SIZE

	for count = 0; count <= perfectIterations; count++ {
		payload := make([]byte, 12)
		byteOffset = count * BLOCK_SIZE
		var length uint32 = BLOCK_SIZE
		if count == perfectIterations {
			length = lastIterationLength
		}
		binary.BigEndian.PutUint32(payload[0:4], uint32(pieceIndex))
		binary.BigEndian.PutUint32(payload[4:8], uint32(byteOffset))
		binary.BigEndian.PutUint32(payload[8:], length)
		fmt.Printf("payload: %v\n", payload)
		fmt.Printf("Sending request for block at offset %d (length: %d)\n", byteOffset, length)
		conn.Write(CreateMessage(6, payload))
	}

	// Receive block data
	byteOffset = 0
	for count > 0 {
		tempData, err := WaitForPeerMessage(conn, 7)
		if err != nil {
			fmt.Printf("Error receiving block data: %v\n", err)
			return err
		}
		begin := binary.BigEndian.Uint32(tempData[4:8])
		copy(data[begin:], tempData[8:])
		fmt.Printf("--Received block data for offset %d--\n", begin)
		count--
	}

	// Verify piece hash
	sum := sha1.Sum(data)
	fmt.Printf("Piece hash: %x\n", sum)
	fmt.Printf("Torrent hash: %x\n", torrent.Info.Pieces[pieceIndex])
	if string(sum[:]) == string(torrent.Info.Pieces[pieceIndex]) {
		fmt.Println("Piece hash verified successfully, writing data to file")
		err = os.WriteFile(outputFile, data, os.ModePerm)
		if err != nil {
			fmt.Printf("Error writing to output file: %v\n", err)
			return err
		}
		fmt.Printf("Data written to file successfully at %s\n", outputFile)
	} else {
		fmt.Println("Piece hash verification failed")
		// return fmt.Errorf("hash not equal")
	}

	fmt.Println("Piece retrieval successful")
	return nil
}

func WaitForPeerMessage(conn net.Conn, expectedMessageID uint8) ([]byte, error) {
	fmt.Printf("Waiting for message with ID %d\n", expectedMessageID)

	for {
		lengthPrefix := make([]byte, 4)
		_, err := io.ReadFull(conn, lengthPrefix)
		if err != nil {
			fmt.Printf("Error reading length prefix: %v\n", err)
			return nil, err
		}
		msgLen := binary.BigEndian.Uint32(lengthPrefix)

		restMessage := make([]byte, msgLen)
		_, err = io.ReadFull(conn, restMessage)
		if err != nil {
			fmt.Printf("Error reading message: %v\n", err)
			return nil, err
		}

		messageId := restMessage[0]
		if messageId != expectedMessageID {
			fmt.Printf("Expected %d, got %d. Ignoring this message.\n", expectedMessageID, messageId)
			continue
		}

		if msgLen == 0 {
			fmt.Println("Received message with no payload")
			return []byte{}, nil
		}
		fmt.Printf("Received payload with ID %d, length %d\n", messageId, msgLen-1)
		return restMessage[1:], nil
	}
}

func CreateMessage(MessageID uint8, payload []byte) []byte {
	msg := make([]byte, 4+1+len(payload))
	binary.BigEndian.PutUint32(msg, uint32(4+1+len(payload)))
	msg[4] = MessageID
	copy(msg[5:], payload)
	fmt.Printf("Created message with ID %d and payload length %d\n", MessageID, len(payload))
	return msg
}
