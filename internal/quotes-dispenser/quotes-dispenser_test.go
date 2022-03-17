package quotesdispenser

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuotesDispenser_LoadJSON(t *testing.T) {
	tests := []struct {
		name       string
		src        io.Reader
		wantQuotes Quotes
		wantErr    bool
	}{
		{
			name: "loaded 3 entries",
			src: strings.NewReader(`[
				{"quoteText": "Genius is one percent inspiration and ninety-nine percent perspiration.", "quoteAuthor": "Thomas Edison"}, 
				{"quoteText": "You can observe a lot just by watching.","quoteAuthor": "Yogi Berra"}, 
				{"quoteText": "A house divided against itself cannot stand.","quoteAuthor": "Abraham Lincoln"}]`,
			),
			wantQuotes: Quotes{
				{Text: "Genius is one percent inspiration and ninety-nine percent perspiration.", Author: "Thomas Edison"},
				{Text: "You can observe a lot just by watching.", Author: "Yogi Berra"},
				{Text: "A house divided against itself cannot stand.", Author: "Abraham Lincoln"},
			},
		},
		{
			name:       "zero entries parsed successfully",
			src:        strings.NewReader(`[]`),
			wantQuotes: Quotes{},
		},
		{
			name:       "bad json returns an error",
			src:        strings.NewReader(`[bad json]`),
			wantQuotes: nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := New()
			err := d.LoadJSON(tt.src)

			assert.Equal(t, tt.wantQuotes, d.quotes)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestQuotesDispenser_GetRandom(t *testing.T) {
	tests := []struct {
		name      string
		quotes    Quotes
		wantQuote bool
		wantErr   error
	}{
		{
			name: "got random quote",
			quotes: Quotes{
				{Text: "Genius is one percent inspiration and ninety-nine percent perspiration.", Author: "Thomas Edison"},
				{Text: "You can observe a lot just by watching.", Author: "Yogi Berra"},
				{Text: "A house divided against itself cannot stand.", Author: "Abraham Lincoln"},
			},
			wantQuote: true,
			wantErr:   nil,
		},
		{
			name:      "failed to get random quote - no quotes",
			wantQuote: false,
			wantErr:   ErrNoQuotes,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := New()
			d.setQuotes(tt.quotes)

			got, err := d.GetRandom()

			assert.ErrorIs(t, err, tt.wantErr)

			if tt.wantQuote {
				assert.Contains(t, d.quotes, got)
			} else {
				assert.Equal(t, Quote{}, got)
			}
		})
	}
}
