package main

import (
	"crypto/tls"
	"errors"
	"io"
	"log"
	"net"
)

var PAYLOAD_MAX = 65535
var PORT = "5973" // If you look at a phone, these are the keys you'd press for LYRE
var HOST = "skarmory"
var ENDPOINT = HOST + ":" + PORT
var KEY_DIR = ".lyre"

func memset(buffer []byte, c byte, n int) {
	for i := range n {
		buffer[i] = c
	}
}

func handleSession(con net.Conn) {
	defer con.Close()

	buffer := make([]byte, PAYLOAD_MAX)
	for {
		numBytes, err := con.Read(buffer)
		if errors.Is(err, io.EOF) {
			log.Printf("Hit EOF, closing connection with %v gracefully", con.RemoteAddr().String())
			break
		} else if err != nil {
			log.Printf("Error during connection: %v", err.Error())
			return
		} else if numBytes > 0 {
			log.Printf("%s", buffer)
			memset(buffer, 0, min(numBytes, PAYLOAD_MAX))
		}
	}
}

func main() {
	log.SetFlags(0)

	cert, err := tls.LoadX509KeyPair(KEY_DIR+"/lyre.crt", KEY_DIR+"/lyre.key")
	if err != nil {
		log.Fatalf("Could not load lyre.crt/lyre.key pair in "+KEY_DIR+": %v", err.Error())
	}

	tlsConf := tls.Config{Certificates: []tls.Certificate{cert}}
	sock, err := tls.Listen("tcp", ENDPOINT, &tlsConf)
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
