package main

import (
	"crypto/tls"
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
var connectionPool = make(map[net.Conn]struct{})

func cleanupConnection(conn net.Conn) {
	connectionPoolMtx.Lock()
	log.Printf("Removing %v from connection pool", conn.RemoteAddr().String())
	delete(connectionPool, conn)
	connectionPoolMtx.Unlock()
	conn.Close()
}

func handleSession(conn net.Conn) {
	connectionPoolMtx.Lock()
	connectionPool[conn] = struct{}{}
	connectionPoolMtx.Unlock()
	defer cleanupConnection(conn)

	msgChannel := make(chan []byte)
	go lyrecom.ListenForMessages(conn, msgChannel)

	for {
		message := <-msgChannel
		log.Printf("[%s]: %s", conn.RemoteAddr().String(), message)
		for outConn := range connectionPool {
			if outConn != conn {
				_, err := outConn.Write(message)
				if err != nil {
					log.Printf("Could not send message from %v to %v", conn.RemoteAddr().String(), outConn.RemoteAddr().String())
				}
			}
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
