package server

import (
	"strconv"
	"sync"
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
	CheckAck(userID uint32, messageID string) bool
	ConfirmAck(userID uint32, messageID string)
	AddAck(userID uint32, messageID string)
}

type AckHandler struct {
	data  map[AckData]bool
	mutex sync.RWMutex
}

func NewAckHandler() IAckHandler {
	return &AckHandler{
		data:  make(map[AckData]bool),
		mutex: sync.RWMutex{},
	}
}

func (handler *AckHandler) CheckAck(to uint32, messageID string) bool {
	handler.mutex.RLock()
	defer handler.mutex.RUnlock()
	return handler.data[Ack{to, messageID}.Value()]
}

func (handler *AckHandler) ConfirmAck(to uint32, messageID string) {
	handler.mutex.Lock()
	defer handler.mutex.Unlock()
	delete(handler.data, Ack{to, messageID}.Value())
}

func (handler *AckHandler) AddAck(to uint32, messageID string) {
	handler.mutex.Lock()
	defer handler.mutex.Unlock()
	handler.data[Ack{to, messageID}.Value()] = true
}
