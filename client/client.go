package main

import (
	"bufio"
	"log"
	"net"
	"os"
)

func main() {
	con, err := net.Dial("tcp", "localhost:5973")
	if err != nil {
		log.Fatalf("Could not connect to server: %v", err.Error())
	}

	defer con.Close()

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
