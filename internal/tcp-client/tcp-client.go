package tcpclient

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/sanyokbig/word-of-wisdom/internal/message"
	"github.com/sanyokbig/word-of-wisdom/internal/methods/simple"
	"github.com/sanyokbig/word-of-wisdom/internal/solver"
	"github.com/sanyokbig/word-of-wisdom/internal/wire"
)

type Wire interface {
	Send(message.MsgType, json.Marshaler) error
	Scanner() *bufio.Scanner
}

type WordsOfWisdom struct {
	Text, Author string
}

type Method interface {
	F(uint64) uint64
}

type TCPClient struct {
	wire Wire

	wordsOfWisdomCh chan *WordsOfWisdom
}

func New(conn net.Conn) *TCPClient {
	return &TCPClient{
		wire:            wire.New(conn),
		wordsOfWisdomCh: make(chan *WordsOfWisdom, 1),
	}
}

func (c *TCPClient) RequestWordsOfWisdom() (*WordsOfWisdom, error) {
	log.Printf("requesting words of wisdom")

	err := c.wire.Send(message.WordsOfWisdomRequest, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	log.Printf("words of wisdom requested")

	wordsOfWisdom, ok := <-c.wordsOfWisdomCh
	if !ok {
		return nil, fmt.Errorf("unable to receive words of wisdom")
	}

	return wordsOfWisdom, nil
}

func (c *TCPClient) Process() {
	scanner := c.wire.Scanner()

	for scanner.Scan() {
		data := scanner.Bytes()

		copied := make([]byte, len(data))
		copy(copied, data)

		c.handleMsg(copied)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("connection closed: %v", err)

		return
	}
}

func (c *TCPClient) handleMsg(data []byte) {
	msg := message.Message{}
	err := msg.UnmarshalJSON(data)
	if err != nil {
		log.Printf("failed to unmarshal msg %s: %v", data, err)

		return
	}

	log.Printf("received server message: %+v", msg)
	switch msg.Type {
	case message.ChallengeRequest:
		c.handleChallengeRequest(msg.Payload)
	case message.WordsOfWisdomResponse:
		c.handleWordsOfWisdomResponse(msg.Payload)
	}
}

func (c *TCPClient) handleChallengeRequest(payload []byte) {
	log.Printf("processing challenge request: %s", payload)

	msg := message.ChallengeRequestPayload{}
	err := msg.UnmarshalJSON(payload)
	if err != nil {
		log.Printf("failed to unmarshal challenge request payload: %v", err)

		return
	}

	now := time.Now().UTC()

	solution, err := c.solve(msg.Xk, msg.N, msg.K, msg.Checksum)
	if err != nil {
		log.Printf("failed to solve challenge: %v", err)

		close(c.wordsOfWisdomCh)

		return
	}

	log.Printf("solution found in %v: %v", time.Since(now), solution)

	err = c.wire.Send(
		message.ChallengeResponse,
		message.ChallengeResponsePayload{
			Y0: solution,
		},
	)
	if err != nil {
		log.Printf("failed to send solution: %v", err)

		return
	}

	log.Printf("solution sent to server")
}

func (c *TCPClient) solve(xk uint64, n, k int, checksum string) (uint64, error) {
	preparedSolver := solver.New(xk, n, k, checksum, simple.New(n))

	solution, ok := preparedSolver.Solve()
	if !ok {
		return 0, fmt.Errorf("solution not found")
	}

	return solution, nil
}

func (c *TCPClient) handleWordsOfWisdomResponse(payload []byte) {
	log.Printf("processing words of wisdom response: %s", payload)

	msg := message.WordsOfWisdomResponsePayload{}
	err := msg.UnmarshalJSON(payload)
	if err != nil {
		log.Printf("failed to unmarshal words of wisdom response payload: %v", err)

		return
	}

	c.wordsOfWisdomCh <- &WordsOfWisdom{
		Text:   msg.Text,
		Author: msg.Author,
	}
}
