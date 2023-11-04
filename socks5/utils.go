package socks5

import (
	"io"
	"net"
)

func readBytes(conn net.Conn, n int) []byte {
	buffer := make([]byte, n)
	if _, err := io.ReadFull(conn, buffer); err != nil {
		panic(err)
	}
	return buffer
}

func writeBytes(conn net.Conn, buffer []byte) {
	if _, err := conn.Write(buffer); err != nil {
		panic(err)
	}
}
