package entity

import "encoding/json"

const (
	MessageTypeUnknown   int = iota
	MessageTypeChallenge     = iota
	MessageTypeResponse
	MessageTypeQuote
)

type Message struct {
	Type int
	Data json.RawMessage
}

type MessageSend struct {
	Type int
	Data any
}

type QuoteMessage struct {
	Quote string
}
