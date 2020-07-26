package sdk

import (
	"github.com/nsqio/go-nsq"
)

type Handler struct {
	msgsChan chan *Message
}

func NewHandler(flight int) *Handler {
	this := &Handler{}
	this.msgsChan = make(chan *Message, flight)

	return this
}

func (h *Handler) HandleMessage(m *nsq.Message) error {
	msg := &Message{
		ID:        MessageIdToStr([nsq.MsgIDLength]byte(m.ID)),
		Body:      m.Body,
		Timestamp: m.Timestamp,
		Attempts:  m.Attempts,
	}

	h.msgsChan <- msg

	return nil
}

func (h *Handler) GetMessage() <-chan *Message {
	return h.msgsChan
}
