package main

import (
	"log"
	"net"

	tcpclient "github.com/sanyokbig/words-of-wisdom/internal/tcp-client"
)

func main() {
	log.Print("starting client...")

	config, err := parse()
	if err != nil {
		log.Printf("failed to parse config: %v", err)

		return
	}

	conn, err := net.Dial("tcp", config.ServerAddr)
	if err != nil {
		log.Printf("failed to connect with server: %v", err)

		return
	}

	tcpClient := tcpclient.New(conn)
	go tcpClient.Process()

	wordsOfWisdom, err := tcpClient.RequestWordsOfWisdom()
	if err != nil {
		log.Printf("failed to receive words of wisdom: %v", err)

		return
	}

	log.Printf("the words of wisdom: \n\n\t \"%v\", %v\n", wordsOfWisdom.Text, wordsOfWisdom.Author)
}
