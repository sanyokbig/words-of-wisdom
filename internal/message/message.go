package message

import "encoding/json"

type Type int

const (
	WordsOfWisdomRequest Type = iota + 1
	WordsOfWisdomResponse
	ChallengeRequest
	ChallengeResponse
)

//easyjson:json
type Message struct {
	Type    Type            `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

//easyjson:json
type ChallengeRequestPayload struct {
	Xk       uint64 `json:"xk"`
	K        int    `json:"k"`
	N        int    `json:"n"`
	Checksum string `json:"checksum"`
}

//easyjson:json
type ChallengeResponsePayload struct {
	Y0 uint64 `json:"y0"`
}

//easyjson:json
type WordsOfWisdomResponsePayload struct {
	Text   string `json:"text"`
	Author string `json:"author"`
}
