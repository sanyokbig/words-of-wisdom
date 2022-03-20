package tcpclient

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/sanyokbig/word-of-wisdom/internal/message"
	"github.com/sanyokbig/word-of-wisdom/internal/wire"
)

type Wire interface {
	Send(message.Type, json.Marshaler) error
	Scanner() *bufio.Scanner
}

type TCPClient struct {
	wire Wire
}

func New(conn net.Conn) *TCPClient {
	return &TCPClient{
		wire: wire.New(conn),
	}
}

func (c *TCPClient) RequestWordsOfWisdom() error {
	log.Printf("requesting words of wisdom")

	err := c.wire.Send(message.WordsOfWisdomRequest, nil)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	log.Printf("words of wisdom requested")

	return nil
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
		log.Printf("connection closed: %v: err", err)

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

	solution, err := c.solve(msg.Xk, msg.Checksum)
	if err != nil {
		log.Printf("failed to solve challenge: %v", err)

		return
	}

	log.Printf("solution found: %v", solution)

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

	return
}

func (c *TCPClient) solve(xk uint64, checksum string) (uint64, error) {
	panic("implement me")
}

func (c *TCPClient) handleWordsOfWisdomResponse(payload []byte) {
	log.Printf("processing words of wisdom response: %s", payload)

	msg := message.WordsOfWisdomResponsePayload{}
	err := msg.UnmarshalJSON(payload)
	if err != nil {
		log.Printf("failed to unmarshal words of wisdom response payload: %v", err)

		return
	}

	log.Printf("the words of wisdom: \n\t \"%v\", %v", msg.Text, msg.Text)

	return
}
