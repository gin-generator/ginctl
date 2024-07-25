package websocket

import (
	"errors"
	"fmt"
	"github.com/gin-generator/ginctl/package/get"
	"github.com/gin-generator/ginctl/package/logger"
	"sync"
	"time"
)

const (
	Max = 1000
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

func NewClientManager() *ClientManager {
	limit := get.Uint("app.max_pool", Max)
	return &ClientManager{
		Register:  make(chan *Client, limit),
		Unset:     make(chan *Client, limit),
		Total:     0,
		Max:       limit,
		Broadcast: make(chan []byte, limit),
		Errs:      make(chan error, limit),
	}
}

// Scheduler Start the websocket scheduler
func (m *ClientManager) Scheduler() {
	for {
		select {
		case client := <-m.Register:
			m.RegisterClient(client)
		case client := <-m.Unset:
			m.UnsetClient(client)
		case message := <-m.Broadcast:
			clients := m.GetAllClient()
			for _, client := range clients {
				client.Send <- message
			}
		case err := <-m.Errs:
			logger.ErrorString("ClientManager", "Start", err.Error())
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
	m.Pool.Range(func(key, value any) bool {
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
	m.Total += 1
}

// UnsetClient Unset client
func (m *ClientManager) UnsetClient(client *Client) {
	err := client.Socket.Close()
	if err != nil {
		m.Errs <- err
	}
	close(client.Send)
	//pubSubs := client.GetAllChan()
	//for _, sub := range pubSubs {
	//	err =
	//}
	m.Pool.Delete(client.Fd)
	m.Total -= 1
}

// Heartbeat The scheduled task clears timeout links
func (m *ClientManager) Heartbeat() {
	EventListener(time.Microsecond*500, func() {
		clients := m.GetAllClient()
		for _, client := range clients {
			if client.IsHeartbeatTimeout(time.Now().Unix()) {
				m.Unset <- client
				fmt.Println(fmt.Sprintf("超时链接 fd: %s 被清理", client.Fd))
			}
		}
	})
}
