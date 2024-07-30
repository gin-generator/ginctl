package websocket

import (
	"errors"
	"sync"
)

type DistributeHandler interface {
	Distribute(client *Client, message []byte) (err error)
}

type Request struct {
	Client *Client
	Send   []byte
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
