package server

import (
	"strconv"
	"sync"
	"time"
)

type Ack struct {
	To        uint32
	MessageID string
}

type AckData string

func (data Ack) Value() AckData {
	return AckData(data.MessageID + "_" + strconv.Itoa(int(data.To)))
}

type IAckHandler interface {
	AssignAckChan(Ack, chan bool)
	WaitForAck(ack Ack, timeout time.Duration)
	ConfirmAck(Ack)
}

type AckHandler struct {
	data  map[AckData]chan bool
	mutex sync.Mutex
}

func NewAckHandler() IAckHandler {
	return &AckHandler{
		data:  make(map[AckData]chan bool),
		mutex: sync.Mutex{},
	}
}

func (handler *AckHandler) AssignAckChan(ack Ack, ch chan bool) {
	handler.mutex.Lock()
	defer handler.mutex.Unlock()
	handler.data[ack.Value()] = ch
}

func (handler *AckHandler) WaitForAck(ack Ack, timeout time.Duration) {
	select {
	case <-time.After(timeout):
		handler.data[ack.Value()] <- false
		return
	}
}

func (handler *AckHandler) ConfirmAck(ack Ack) {
	value := ack.Value()
	handler.mutex.Lock()
	defer handler.mutex.Unlock()
	if ackChan, ok := handler.data[value]; ok {
		ackChan <- true
		delete(handler.data, value)
	}
}
