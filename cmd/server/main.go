package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/sanyokbig/words-of-wisdom/internal/challenger"
	cryptorand "github.com/sanyokbig/words-of-wisdom/internal/crypto-rand"
	quotesdispenser "github.com/sanyokbig/words-of-wisdom/internal/quotes-dispenser"
	tcpserver "github.com/sanyokbig/words-of-wisdom/internal/tcp-server"
)

func main() {
	log.Print("starting server...")

	config, err := parse()
	if err != nil {
		log.Panicf("failed to parse config: %v", err)
	}

	// Load quotes file
	quotesFile, err := os.Open(config.QuotesSource)
	if err != nil {
		log.Printf("failed to open quotes file: %v", err)

		return
	}

	// Prepare quotes dispenser
	quotesDispenser := quotesdispenser.New()
	err = quotesDispenser.LoadJSON(quotesFile)
	if err != nil {
		log.Printf("failed to load json: %v", err)

		return
	}

	err = quotesFile.Close()
	if err != nil {
		log.Printf("failed to close quotes file: %v", err)

		return
	}

	listenAddr := ":" + strconv.Itoa(config.TCPPort)

	tcpServer := tcpserver.New(quotesDispenser, challenger.New(cryptorand.Uint64))
	go func() {
		err := tcpServer.ListenAndServe(listenAddr)
		if err != nil {
			log.Printf("failed to listen and serve")
		}
	}()

	log.Printf("listeninig on %v", listenAddr)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

	sig := <-ch
	log.Printf("got signal %v", sig)
}
