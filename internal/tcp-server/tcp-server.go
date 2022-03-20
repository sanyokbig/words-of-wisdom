package tcpserver

import (
	"fmt"
	"log"
	"net"

	"github.com/sanyokbig/word-of-wisdom/internal/challenger"
	"github.com/sanyokbig/word-of-wisdom/internal/methods/simple"
)

type WordsOfWisdom interface {
	Get() (string, string, error)
}

type Challenger interface {
	Prepare(method challenger.Method, n, k int) *challenger.Challenge
}

type TCPServer struct {
	listener net.Listener

	wordsOfWisdom WordsOfWisdom
	challenger    Challenger
}

func New(wordsOfWisdom WordsOfWisdom, ch Challenger) *TCPServer {
	return &TCPServer{
		wordsOfWisdom: wordsOfWisdom,
		challenger:    ch,
	}
}

func (s *TCPServer) ListenAndServe(listenAddr string) error {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return fmt.Errorf("failed to listen on %v: %w", listenAddr, err)
	}

	s.listener = listener

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v", err)

			continue
		}

		go s.handleConn(conn)
	}
}

func (s *TCPServer) prepareChallenge() *challenger.Challenge {
	// Right now both bitSize and depth parameters are constant, but can be dynamically adjusted based on a number of
	// active clients or other metrics and conditions if needed
	var (
		// A bit size mostly affects only client and slightly server
		bitSize = 21

		// Depth affects both server and client execution time
		depth = 64
	)

	return s.challenger.Prepare(simple.New(bitSize), bitSize, depth)
}

func (s *TCPServer) handleConn(conn net.Conn) {
	client := NewClient(conn, s.wordsOfWisdom, s)

	client.Process()
}
