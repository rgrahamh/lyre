package main

import (
	"bufio"
	"crypto/tls"
	"log"
	"lyrecom"
	"os"
)

var PAYLOAD_MAX = 65535

func inputReader(con *tls.Conn) {
	reader := bufio.NewReader(os.Stdin)

	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Could not read string from stdin")
		}

		// We don't want to send a newline back at the end
		buff := ([]byte)(text)
		buff[len(buff)-1] = 0

		_, err = con.Write(buff)
		if err != nil {
			log.Fatalf("Could not write to server: %v", err.Error())
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

	tlsConf := tls.Config{InsecureSkipVerify: true}
	con, err := tls.Dial("tcp", "trashsuite.games:5973", &tlsConf)
	if err != nil {
		log.Fatalf("Could not connect to server: %v", err.Error())
	}
	defer con.Close()

	messageChannel := make(chan []byte)
	go lyrecom.ListenForMessages(con, messageChannel)
	go handleMessages(messageChannel)
	go inputReader(con)
}
