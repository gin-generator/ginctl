package websocket

import (
	"errors"
	"fmt"
	"github.com/gin-generator/ginctl/package/get"
	"github.com/gin-generator/ginctl/package/logger"
	rds "github.com/gin-generator/ginctl/package/redis"
	t "github.com/gin-generator/ginctl/package/time"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"sync"
	"time"
)

type Client struct {
	Fd            string          // 每个连接唯一标识
	Addr          string          // 客户端ip地址
	Socket        *websocket.Conn // 用户连接
	Send          chan []byte     // 待发送的数据
	FirstTime     int64           // 首次连接事件
	HeartbeatTime int64           // 用户上次心跳时间
	Timeout       int64           // 断连时间
	Channel       sync.Map        // 订阅频道
}

func NewClient(addr string, socket *websocket.Conn) *Client {
	First := time.Now().Unix()
	limit := get.Uint("app.max_pool", Max)
	return &Client{
		Fd:            uuid.NewV4().String(),
		Addr:          addr,
		Socket:        socket,
		Send:          make(chan []byte, limit),
		FirstTime:     First,
		HeartbeatTime: First,
	}
}

func (c *Client) Close() {
	err := c.Socket.Close()
	if err != nil {
		logger.ErrorString("Websocket", "Close", err.Error())
		return
	}
	close(c.Send)
}

// Read client data
func (c *Client) Read() {

	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			if closeErr, ok := err.(*websocket.CloseError); ok {
				logger.ErrorString("Websocket", "Read",
					fmt.Sprintf("%s close code: %v\n", t.Time{}.Local(), closeErr.Code))
			}
			return
		}

		err = Distribute(c, message)
		if err != nil {
			logger.ErrorString("Read", "Distribute",
				fmt.Sprintf("Message distribution error: %s, fd: %s", err.Error(), c.Fd))
			return
		}
	}

}

// Write Send data to the client
func (c *Client) Write() {

	for bytes := range c.Send {
		if err := c.Socket.WriteMessage(websocket.TextMessage, bytes); err != nil {
			logger.ErrorString("Websocket", "Write",
				fmt.Sprintf("send message err: %s, address: %s,fd: %s", err.Error(), c.Addr, c.Fd))
		}
	}

}

// SendMessage Message distribution
func (c *Client) SendMessage(message []byte) {

	if c == nil {
		logger.ErrorString("Websocket", "SendMessage",
			fmt.Sprintf("Not found websocket client,fd: %s", c.Fd))
		return
	}

	c.Send <- message
}

// CreatChan Creat channel
func (c *Client) CreatChan() (channel string, err error) {
	channel = uuid.NewV4().String()
	// TODO: The publish-subscribe model will be implemented later here using other drivers
	pubSub := rds.Rds.Subscribe(channel)
	_, err = pubSub.Receive(rds.Rds.Context)
	if err != nil {
		return
	}
	c.Channel.LoadOrStore(channel, pubSub)
	return
}

// GetChan Get channel
func (c *Client) GetChan(channel string) (pubSub *redis.PubSub, err error) {
	// TODO: The publish-subscribe model will be implemented later here using other drivers
	value, ok := c.Channel.Load(channel)
	if !ok {
		return nil, errors.New("not found channel")
	}
	pubSub, ok = value.(*redis.PubSub)
	if !ok {
		return nil, errors.New("channel type error")
	}
	return
}

// GetAllChan Get all channels
func (c *Client) GetAllChan() (pubSubs []*redis.PubSub) {
	c.Channel.Range(func(key, value any) bool {
		pubSub, ok := value.(*redis.PubSub)
		if ok {
			pubSubs = append(pubSubs, pubSub)
		}
		return true
	})
	return
}

// Publish a message
func (c *Client) Publish(channel string, message []byte) (err error) {
	err = rds.Rds.Publish(channel, string(message))
	return
}

// Subscribe to messages for long links
func (c *Client) Subscribe(channel string) (err error) {
	// TODO: The publish-subscribe model will be implemented later here using other drivers
	pubSub := rds.Rds.Subscribe(channel)
	_, err = pubSub.Receive(rds.Rds.Context)
	if err != nil {
		return err
	}
	c.Channel.LoadOrStore(channel, pubSub)
	return
}

// Receive subscription messages
func (c *Client) Receive() {
	var pubSubs []*redis.PubSub
	EventListener(time.Millisecond*1000, func() {
		pubSubs = c.GetAllChan()
	})
	for _, sub := range pubSubs {
		go func(sub *redis.PubSub) {
			ch := sub.Channel()
			for message := range ch {
				c.Send <- []byte(message.Payload)
			}
		}(sub)
	}
}

// Unsubscribe unsubscribe
func (c *Client) Unsubscribe(channel string) (err error) {
	pubSub, err := c.GetChan(channel)
	if err != nil {
		return
	}
	err = pubSub.Unsubscribe(rds.Rds.Context, channel)
	if err != nil {
		return
	}
	return
}

func (c *Client) SetHeartbeatTime(currentTime int64) {
	c.HeartbeatTime = currentTime
}

func (c *Client) IsHeartbeatTimeout(currentTime int64) (timeout bool) {
	if c.HeartbeatTime+c.Timeout <= currentTime {
		timeout = true
	}
	return
}