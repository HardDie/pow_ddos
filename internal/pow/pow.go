package pow

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/HardDie/pow_ddos/internal/logger"
	"github.com/HardDie/pow_ddos/internal/utils"
)

type Block struct {
	Timestamp  time.Time
	Data       []byte
	Difficulty int
	Nonce      int64
	Hash       []byte
}

func NewBlock(size, difficulty int) (*Block, error) {
	buf, err := utils.RandomBytes(size)
	if err != nil {
		return nil, err
	}

	return &Block{
		Timestamp:  time.Now(),
		Data:       buf,
		Difficulty: difficulty,
	}, nil
}

func (b Block) CalcHash() []byte {
	record := strconv.FormatInt(b.Timestamp.Unix(), 10) +
		string(b.Data) +
		strconv.FormatInt(b.Nonce, 10)
	hash := sha256.Sum256([]byte(record))
	return hash[:]
}

func (b *Block) POW(threads int) error {
	if threads < 1 {
		threads = 1
	}
	logger.Debug.Println("Run threads:", threads)

	var resNonce int64
	ctxIsFound, ctxCancel := context.WithCancel(context.Background())
	isDone := make(chan struct{})

	guessRoutine := func(id int, from, to int64) {
		routineIsDone := false
		newBlock := *b
		for nonce := from; nonce < to && !routineIsDone; nonce++ {
			// Calculate hash
			newBlock.Nonce = nonce
			hash := newBlock.CalcHash()

			// Check if result was found
			if utils.CountLeadingZeros(hash) >= newBlock.Difficulty {
				// Notify other goroutines result is found
				ctxCancel()
				// Even if we found more than one answer, it doesn't matter, we can use any of them
				resNonce = nonce
				break
			}

			// Check if other goroutine found result
			select {
			case <-ctxIsFound.Done():
				routineIsDone = true
			default:
			}
		}
		isDone <- struct{}{}
	}

	timeStart := time.Now()

	// Run all goroutines
	guessStep := int64(math.MaxInt64 / threads)
	from := int64(0)
	to := guessStep
	for i := 0; i < threads; i++ {
		go guessRoutine(i, from, to)
		from, to = to+1, to+guessStep
		if (math.MaxInt64 - to) < guessStep {
			to = math.MaxInt64
		}
	}

	// Wait until all goroutines is done
	var count int
	for {
		select {
		case <-isDone:
			count++
		}
		if count == threads {
			break
		}
	}

	timeEnd := time.Now()
	logger.Debug.Println("Guess time:", timeEnd.Sub(timeStart))

	if resNonce == 0 {
		return fmt.Errorf("hash was not found")
	}

	b.Nonce = resNonce
	b.Hash = b.CalcHash()
	return nil
}
