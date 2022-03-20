package tcpserver

import (
	"fmt"
	"log"
	"net"
)

type TCPServer struct {
	listener net.Listener

	wordsOfWisdom WordsOfWisdom
}

func New(wordsOfWisdom WordsOfWisdom) *TCPServer {
	return &TCPServer{
		wordsOfWisdom: wordsOfWisdom,
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

func (s *TCPServer) Close() error {
	return s.listener.Close()
}

func (s *TCPServer) handleConn(conn net.Conn) {
	NewClient(conn, s.wordsOfWisdom).Process()
}
