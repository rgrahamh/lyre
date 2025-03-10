package main

import (
	"errors"
	"io"
	"log"
	"net"
)

var PAYLOAD_MAX = 65535
var PORT = "5973" // If you look at a phone, these are the keys you'd press for LYRE
var ADDRESS = "0.0.0.0"
var ENDPOINT = ADDRESS + ":" + PORT

func handleSession(con net.Conn) {
	defer con.Close()

	buffer := make([]byte, PAYLOAD_MAX)
	for {
		numBytes, err := con.Read(buffer)
		if errors.Is(err, io.EOF) {
			log.Printf("Hit EOF, closing connection with %v gracefully", con.RemoteAddr().String())
			break
		} else if err != nil {
			log.Fatalf("Error during connection: %v", err.Error())
		} else if numBytes > 0 {
			log.Printf("%s", buffer)
		}
	}
}

func main() {
	log.SetFlags(0)

	sock, err := net.Listen("tcp", ENDPOINT)
	if err != nil {
		log.Fatalf("Could not open a server on port %v: %v", ENDPOINT, err.Error())
	}

	defer sock.Close()

	for {
		con, err := sock.Accept()
		if err != nil {
			log.Fatalf("Error receiving connection: %v", err.Error())
		}

		go handleSession(con)
	}

}
