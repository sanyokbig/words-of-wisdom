package tcpserver

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/sanyokbig/words-of-wisdom/internal/challenger"
	"github.com/sanyokbig/words-of-wisdom/internal/message"
	"github.com/sanyokbig/words-of-wisdom/internal/wire"
)

type Wire interface {
	Send(message.MsgType, json.Marshaler) error
	Scanner() *bufio.Scanner
}

type Client struct {
	wire Wire

	wordsOfWisdom WordsOfWisdom

	challengePreparer challengePreparer
	challenge         *challenger.Challenge
}

// challengePreparer is a simplified version of Challenger, so that we can receive a challenge with a different
// difficulty depending on external factors like server load that client should not worry and know about.
type challengePreparer interface {
	prepareChallenge() *challenger.Challenge
}

func NewClient(conn net.Conn, wordsOfWisdom WordsOfWisdom, preparer challengePreparer) *Client {
	return &Client{
		wire:              wire.New(conn),
		wordsOfWisdom:     wordsOfWisdom,
		challengePreparer: preparer,
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
		log.Printf("scanner err: %v: err", err)
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
		c.handleWordsOfWisdomRequest()
	case message.ChallengeResponse:
		c.handleChallengeResponse(msg.Payload)
	}
}

func (c *Client) handleWordsOfWisdomRequest() {
	log.Printf("processing words of wisdom request")

	if c.challenge != nil {
		log.Printf("client already received challenge")

		return
	}

	now := time.Now().UTC()
	challenge := c.challengePreparer.prepareChallenge()
	log.Printf("challenged prepared in %v", time.Since(now))

	err := c.wire.Send(
		message.ChallengeRequest,
		message.ChallengeRequestPayload{
			Xk:       challenge.Xk,
			K:        challenge.K,
			N:        challenge.N,
			Checksum: challenge.Checksum,
		})
	if err != nil {
		log.Printf("failed to send challenge to client: %v", err)

		return
	}

	c.challenge = challenge
	log.Printf("client is challenged with %+v", challenge)
}

func (c *Client) handleChallengeResponse(payload []byte) {
	log.Printf("processing challenge response: %s", payload)

	if c.challenge == nil {
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

	log.Printf("solution is valid, granting a words of wisdom")

	err = c.grantWordsOfWisdom()
	if err != nil {
		log.Printf("failed to grant a words of wisdom: %v", err)
	}
}

func (c *Client) validateSolution(y0 uint64) bool {
	if c.challenge == nil {
		return false
	}

	return c.challenge.X0 == y0
}

func (c *Client) grantWordsOfWisdom() error {
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

	c.challenge = nil

	log.Printf("words of wisdom granted")

	return nil
}
