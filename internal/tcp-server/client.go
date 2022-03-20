package tcpserver

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

type Client struct {
	wire Wire

	challengeInfo *challengeInfo

	wordsOfWisdom WordsOfWisdom
}

type WordsOfWisdom interface {
	Get() (string, string, error)
}

func NewClient(conn net.Conn, wordsOfWisdom WordsOfWisdom) *Client {
	return &Client{
		wire:          wire.New(conn),
		wordsOfWisdom: wordsOfWisdom,
	}
}

func (c *Client) Process() {
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

func (c *Client) handleMsg(data []byte) {
	msg := message.Message{}
	err := msg.UnmarshalJSON(data)
	if err != nil {
		log.Printf("failed to unmarshal msg %s: %v", data, err)

		return
	}

	log.Printf("received client message: %+v", msg)
	switch msg.Type {
	case message.WordsOfWisdomRequest:
		c.handleWordOfWisdomRequest()
	case message.ChallengeResponse:
		c.handleChallengeResponse(msg.Payload)
	}
}

func (c *Client) handleWordOfWisdomRequest() {
	log.Printf("processing word of wisdom request")

	if c.challengeInfo != nil {
		log.Printf("client already received challenge")

		return
	}

	ci := makeChallenge()

	err := c.wire.Send(
		message.ChallengeRequest,
		message.ChallengeRequestPayload{
			Xk:       c.challengeInfo.xk,
			Checksum: c.challengeInfo.checksum,
		})
	if err != nil {
		log.Printf("failed to send challenge to client: %v", err)

		return
	}

	c.challengeInfo = ci
	log.Printf("client is challenged with %+v", ci)
}

type challengeInfo struct {
	x0, xk   uint64
	checksum string
}

func makeChallenge() *challengeInfo {
	panic("implement me")
}

func (c *Client) handleChallengeResponse(payload []byte) {
	log.Printf("processing challenge response: %s", payload)

	if c.challengeInfo == nil {
		log.Printf("client is not challenged")

		return
	}

	msg := message.ChallengeResponsePayload{}
	err := msg.UnmarshalJSON(payload)
	if err != nil {
		log.Printf("failed to unmarshal challenge response payload: %v", err)

		return
	}

	if !c.validateSolution(msg.Y0) {
		log.Printf("client solution is wrong")

		return
	}

	log.Printf("solution is valid, granting a word of wisdom")

	err = c.grandWordOfWisdom()
	if err != nil {
		log.Printf("failed to grant a word of wisdom: %v", err)

		return
	}
}

func (c *Client) validateSolution(y0 uint64) bool {
	panic("implement me")
}

func (c *Client) grandWordOfWisdom() error {
	text, author, err := c.wordsOfWisdom.Get()
	if err != nil {
		return fmt.Errorf("failed to get: %w", err)
	}

	payload := message.WordsOfWisdomResponsePayload{
		Text:   text,
		Author: author,
	}

	err = c.wire.Send(message.WordsOfWisdomResponse, payload)
	if err != nil {
		return fmt.Errorf("failed to write msg: %w", err)
	}

	c.challengeInfo = nil

	log.Printf("word of wisdom granted")

	return nil
}
