package websocket

import (
	"encoding/json"
	"errors"
	"sync"
)

type Request struct {
	Client *Client
	Send   []byte
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type Message struct {
	Event   string `json:"event"`
	Request string `json:"request,omitempty"`
}

type DisposeFunc func(request *Request) (response *Response)

var (
	handlers sync.Map
)

// Register event to memory
func Register(event string, handler DisposeFunc) {
	handlers.LoadOrStore(event, handler)
}

// GetHandler Get register's func
func GetHandler(event string) (handler DisposeFunc, err error) {
	value, ok := handlers.Load(event)
	if !ok {
		return nil, errors.New("current event not supported")
	}
	handler, okk := value.(DisposeFunc)
	if !okk {
		return nil, errors.New("handler type error")
	}
	return
}

// Distribute Data distribution processing
func Distribute(client *Client, message []byte) (err error) {
	r := &Message{}
	err = json.Unmarshal(message, r)
	if err != nil {
		return
	}

	handler, err := GetHandler(r.Event)
	if err != nil {
		return
	}

	response := handler(&Request{
		Client: client,
		Send:   []byte(r.Request),
	})

	bytes, err := json.Marshal(response)
	if err != nil {
		return
	}

	client.SendMessage(bytes)
	return
}
