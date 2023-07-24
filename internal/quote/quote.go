package quote

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"

	"github.com/HardDie/pow_ddos/internal/config"
)

type Quote struct {
	quotes []string
}

func New(config config.Quote) (*Quote, error) {
	file, err := os.Open(config.Path)
	if err != nil {
		return nil, fmt.Errorf("error open file with quotes: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading quotes: %w", err)
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("quotes file is empty")
	}

	return &Quote{
		quotes: strings.Split(strings.TrimSpace(string(data)), "\n"),
	}, nil
}

func (q *Quote) GetRandomQuote() string {
	if len(q.quotes) == 0 {
		return ""
	}
	return q.quotes[rand.Intn(len(q.quotes))]
}
