package main

import (
	"bufio"
	"crypto/tls"
	"errors"
	"io"
	"log"
	"lyrecom"
	"os"
	"os/signal"
	"syscall"
)

func inputReader(con *tls.Conn) {
	reader := bufio.NewReader(os.Stdin)

	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Could not read string from stdin")
		}

		// We'll remove the last byte, to strip the newline
		buff := ([]byte)(text[0 : len(text)-1])

		_, err = con.Write(buff)
		if err != nil {
			log.Fatalf("Could not write to server: %v", err.Error())
		}
	}
}

func ListenForMessages(con *tls.Conn, msgChannel chan []byte) {
	buffer := make([]byte, lyrecom.PAYLOAD_MAX)
	for {
		numBytes, err := (*con).Read(buffer)
		if errors.Is(err, io.EOF) {
			log.Printf("Hit EOF, closing connection with %v gracefully", (*con).RemoteAddr().String())
			break
		} else if err != nil {
			log.Printf("Error during connection: %v", err.Error())
			break
		} else if numBytes > 0 {
			msgChannel <- buffer[0:numBytes]
		}
	}
}

func handleMessages(messageChannel chan []byte) {
	for {
		log.Printf("%s", <-messageChannel)
	}
}

func main() {
	log.SetFlags(0)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	tlsConf := tls.Config{InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", "trashsuite.games:5973", &tlsConf)
	if err != nil {
		log.Fatalf("Could not connect to server: %v", err.Error())
	}
	defer conn.Close()

	messageChannel := make(chan []byte)
	go ListenForMessages(conn, messageChannel)
	go handleMessages(messageChannel)
	go inputReader(conn)

	<-sig
}
