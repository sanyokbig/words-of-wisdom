package main

import (
	"log"
	"math/rand"
	"net"
	"time"

	tcpclient "github.com/sanyokbig/word-of-wisdom/internal/tcp-client"
)

func main() {
	log.Print("starting client...")
	rand.Seed(time.Now().UnixNano())

	config, err := parse()
	if err != nil {
		log.Panicf("failed to parse config: %v", err)
	}

	conn, err := net.Dial("tcp", config.ServerAddr)
	if err != nil {
		log.Panicf("failed to connect with server: %v", err)
	}

	tcpClient := tcpclient.New(conn)
	go tcpClient.Process()

	wordsOfWisdom, err := tcpClient.RequestWordsOfWisdom()
	if err != nil {
		log.Panicf("failed to request words of wisdom: %v", err)
	}

	log.Printf("the words of wisdom: \n\t \"%v\", %v", wordsOfWisdom.Text, wordsOfWisdom.Text)
}
