package websocket

import (
	"errors"
	"fmt"
	"github.com/gin-generator/ginctl/package/get"
	"github.com/gin-generator/ginctl/package/logger"
	"github.com/go-redis/redis/v8"
	"sync"
)

const (
	Max = 1000
)

var (
	Manager *ClientManager
	once    sync.Once
)

// ClientManager Client pool manager
type ClientManager struct {
	Pool      sync.Map
	Register  chan *Client
	Unset     chan *Client
	Total     uint
	Max       uint
	Broadcast chan []byte
	Errs      chan error
}

func NewClientManager() {
	limit := get.Uint("app.max_pool", Max)
	once.Do(func() {
		Manager = &ClientManager{
			Register:  make(chan *Client, limit),
			Unset:     make(chan *Client, limit),
			Total:     0,
			Max:       limit,
			Broadcast: make(chan []byte, limit),
			Errs:      make(chan error, limit),
		}
	})

	go Manager.Scheduler()
}

// Scheduler Start the websocket scheduler
func (m *ClientManager) Scheduler() {
	for {
		select {
		case client := <-m.Register:
			m.RegisterClient(client)
		case client := <-m.Unset:
			m.Close(client)
		case message := <-m.Broadcast:
			m.Pool.Range(func(_, value any) bool {
				client, ok := value.(*Client)
				if ok {
					client.Send <- message
				}
				return true
			})
		case err := <-m.Errs:
			logger.ErrorString("ClientManager", "Scheduler", err.Error())
		}
	}
}

// GetClient Obtain the client according to fd
func (m *ClientManager) GetClient(fd string) (client *Client, err error) {
	value, ok := m.Pool.Load(fd)
	if !ok {
		return nil, errors.New("no client found")
	}

	client, ok = value.(*Client)
	if !ok {
		return nil, errors.New("client error")
	}
	return
}

// GetAllClient Get all client
func (m *ClientManager) GetAllClient() (clients []*Client) {
	m.Pool.Range(func(_, value any) bool {
		client, ok := value.(*Client)
		if ok {
			clients = append(clients, client)
		}
		return true
	})
	return clients
}

// RegisterClient Register client
func (m *ClientManager) RegisterClient(client *Client) {
	m.Pool.Store(client.Fd, client)
	m.Total++
}

// Close Unset client
func (m *ClientManager) Close(client *Client) {
	err := client.Socket.Close()
	if err != nil {
		m.Errs <- err
	}
	close(client.Send)
	channels, _ := client.GetAllChan()
	for _, channel := range channels {
		err = client.Unsubscribe(channel)
		if err != nil {
			m.Errs <- err
		}
	}
	client.OwnerChannel.Range(func(key, value any) bool {
		pubSub, ok := value.(*redis.PubSub)
		if ok {
			m.Errs <- pubSub.Close()
		}
		return true
	})
	m.Pool.Delete(client.Fd)
	m.Total--
	logger.InfoString("ClientManager", "UnsetClient",
		fmt.Sprintf("websocket timeout, fd: %s be cleared", client.Fd))
}
