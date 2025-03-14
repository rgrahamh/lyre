package lyrecom

import (
	"errors"
	"io"
	"log"
	"net"
)

var PAYLOAD_MAX = 65535

func Memset(buffer []byte, c byte, n int) {
	for i := range n {
		buffer[i] = c
	}
}

func ListenForMessages(con net.Conn, msgChannel chan []byte) {
	for {
		buffer := make([]byte, PAYLOAD_MAX)
		numBytes, err := con.Read(buffer)
		if errors.Is(err, io.EOF) {
			log.Printf("Hit EOF, closing connection with %v gracefully", con.RemoteAddr().String())
			break
		} else if err != nil {
			log.Printf("Error during connection: %v", err.Error())
			break
		} else if numBytes > 0 {
			log.Printf("[%s]: %s", con.RemoteAddr().String(), buffer)
			msgChannel <- buffer
		}
	}
}
