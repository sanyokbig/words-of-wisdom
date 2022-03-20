package wire

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"

	"github.com/sanyokbig/word-of-wisdom/internal/message"
)

type Wire struct {
	conn    net.Conn
	scanner *bufio.Scanner
}

func New(conn net.Conn) *Wire {
	s := bufio.NewScanner(conn)

	s.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	s.Split(scanDocuments)

	return &Wire{
		conn:    conn,
		scanner: s,
	}
}

// Scan until line break is encountered.
func scanDocuments(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.Index(data, []byte("\n")); i >= 0 {
		return i + 1, data[0:i], nil
	}

	if atEOF {
		return len(data), data, nil
	}

	// Request more data.
	return 0, nil, nil
}

func (w *Wire) Send(msgType message.MsgType, payload json.Marshaler) (err error) {
	var rawPayload []byte

	if payload != nil {
		rawPayload, err = payload.MarshalJSON()
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %w", err)
		}
	}

	data, err := message.Message{Type: msgType, Payload: rawPayload}.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	_, err = fmt.Fprintf(w.conn, "%s\n", data)
	if err != nil {
		return fmt.Errorf("failed to write msg: %w", err)
	}

	return nil
}

func (w *Wire) Scanner() *bufio.Scanner {
	return w.scanner
}
