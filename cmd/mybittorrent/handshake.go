package main

import (
	"io"
	"net"
)

func Handshake(peer string, hash *[]byte) (net.Conn, []byte, error) {
	conn, err := net.Dial("tcp", peer)
	if err != nil {
		return conn, nil, err
	}
	var buf []byte
	buf = append(buf, 19)                                // 01 byte
	buf = append(buf, []byte("BitTorrent protocol")...)  // 19 bytes
	buf = append(buf, make([]byte, 8)...)                // 08 bytes
	buf = append(buf, *hash...)              // 20 bytes
	buf = append(buf, []byte("00112233445566778899")...) // 20 bytes
	_, err = conn.Write(buf)
	if err != nil {
		return nil, nil, err
	}
	answer := make([]byte, 68)
	io.ReadFull(conn, answer)
	return conn, answer, nil
}