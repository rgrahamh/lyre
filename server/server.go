package main

import (
	"crypto/tls"
	"errors"
	"io"
	"log"
	"lyrecom"
	"net"
	"sync"
)

var PORT = "5973" // If you look at a phone, these are the keys you'd press for LYRE
var HOST = "skarmory"
var ENDPOINT = HOST + ":" + PORT
var KEY_DIR = ".lyre"

var connectionPoolMtx sync.Mutex
var connectionPool = make(map[*net.Conn]struct{})

func ListenForMessages(conn *net.Conn, msgChannel chan []byte) {
	buffer := make([]byte, lyrecom.PAYLOAD_MAX)
	for {
		numBytes, err := (*conn).Read(buffer)
		if errors.Is(err, io.EOF) {
			log.Printf("Hit EOF, closing connection with %v gracefully", (*conn).RemoteAddr().String())
			break
		} else if err != nil {
			log.Printf("Error during connection: %v", err.Error())
			break
		} else if numBytes > 0 {
			msgChannel <- buffer[0:numBytes]
		}
	}
}

func cleanupConnection(conn *net.Conn) {
	connectionPoolMtx.Lock()
	_, shouldClose := connectionPool[conn]
	connectionPoolMtx.Unlock()

	if shouldClose {
		delete(connectionPool, conn)
		(*conn).Close()
	}
}

func handleSession(conn net.Conn) {
	connectionPoolMtx.Lock()
	connectionPool[&conn] = struct{}{}
	connectionPoolMtx.Unlock()

	defer cleanupConnection(&conn)

	msgChannel := make(chan []byte)
	go ListenForMessages(&conn, msgChannel)

	for {
		message := <-msgChannel
		log.Printf("[%s]: %s", conn.RemoteAddr().String(), message)
		connectionPoolMtx.Lock()
		for outConn := range connectionPool {
			if outConn != &conn {
				_, err := (*outConn).Write(message)
				if err != nil {
					log.Printf("Could not send message from %v to %v; removing from connection pool", conn.RemoteAddr().String(), (*outConn).RemoteAddr().String())
					delete(connectionPool, &conn)
					conn.Close()
				}
			}
		}
		connectionPoolMtx.Unlock()
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
		conn, err := sock.Accept()
		if err != nil {
			log.Fatalf("Error receiving connection: %v", err.Error())
		}

		go handleSession(conn)
	}
}
