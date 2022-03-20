package quotesdispenser

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"sync"
)

var ErrNoQuotes = errors.New("no quotes")

type QuotesDispenser struct {
	mux sync.RWMutex

	quotes    Quotes
	quotesLen int
}

//easyjson:json
type Quote struct {
	Text   string `json:"quoteText"`
	Author string `json:"quoteAuthor"`
}

//easyjson:json
type Quotes []Quote

func New() *QuotesDispenser {
	return &QuotesDispenser{}
}

// LoadJSON will replace existing quotes with new ones.
func (d *QuotesDispenser) LoadJSON(src io.Reader) error {
	raw, err := io.ReadAll(src)
	if err != nil {
		return fmt.Errorf("read all: %w", err)
	}

	var quotes Quotes
	err = quotes.UnmarshalJSON(raw)
	if err != nil {
		return fmt.Errorf("unmarshal json: %w", err)
	}

	d.setQuotes(quotes)

	return nil
}

func (d *QuotesDispenser) setQuotes(quotes Quotes) {
	d.mux.Lock()
	d.quotes = quotes
	d.quotesLen = len(quotes)
	d.mux.Unlock()
}

func (d *QuotesDispenser) Get() (string, string, error) {
	d.mux.RLock()
	defer d.mux.RUnlock()

	if d.quotesLen == 0 {
		return "", "", ErrNoQuotes
	}

	//nolint:gosec
	idx := rand.Intn(d.quotesLen - 1)

	quote := d.quotes[idx]

	return quote.Text, quote.Author, nil
}
