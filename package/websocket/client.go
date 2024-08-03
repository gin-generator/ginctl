package websocket

import (
	"errors"
	"fmt"
	"github.com/gin-generator/ginctl/package/get"
	"github.com/gin-generator/ginctl/package/logger"
	rds "github.com/gin-generator/ginctl/package/redis"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"sync"
	"time"
)

type Client struct {
	Fd            string             // 每个连接唯一标识
	Addr          string             // 客户端ip地址
	Socket        *websocket.Conn    // 用户连接
	Send          chan []byte        // 待发送的数据
	Channel       sync.Map           // 订阅频道
	OwnerChannel  sync.Map           // 自己创建的频道
	receive       chan *redis.PubSub // 订阅通知
	FirstTime     int64              // 首次连接事件
	HeartbeatTime int64              // 用户上次心跳时间
	Timeout       int64              // 超时断连时间
	Protocol      int                // 协议类型
	close         chan struct{}
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
		receive:       make(chan *redis.PubSub, 10),
		Timeout:       get.Int64("app.heartbeat_timeout", 600),
		Protocol:      websocket.TextMessage,
		close:         make(chan struct{}, 1),
	}
}

func (c *Client) Close() {
	c.close <- struct{}{}
	Manager.Unset <- c
}

// Read client data
func (c *Client) Read() {

	var handler DistributeHandler
	for {
		messageType, message, err := c.Socket.ReadMessage()
		if err != nil {
			if closeErr, ok := err.(*websocket.CloseError); ok {
				logger.InfoString("Websocket", "Read",
					fmt.Sprintf("close time: %s, close code: %d",
						time.Now().Format("2006-01-02 15:04:05"), closeErr.Code))
			}
			return
		}

		c.Protocol = messageType
		switch messageType {
		case websocket.TextMessage:
			handler = NewJsonHandler()
		case websocket.BinaryMessage:
			handler = NewProtoHandler()
		}

		err = handler.Distribute(c, message)
		if err != nil {
			logger.ErrorString("Read", "Distribute",
				fmt.Sprintf("%s, fd: %s", err.Error(), c.Fd))
		}
	}

}

// Write Send data to the client
func (c *Client) Write() {

	for bytes := range c.Send {
		if err := c.Socket.WriteMessage(c.Protocol, bytes); err != nil {
			if closeErr, ok := err.(*websocket.CloseError); ok {
				logger.InfoString("Websocket", "Write",
					fmt.Sprintf("close time: %s, close code: %d",
						time.Now().Format("2006-01-02 15:04:05"), closeErr.Code))
				return
			}
			logger.ErrorString("Websocket", "Write",
				fmt.Sprintf("%s, address: %s, fd: %s", err.Error(), c.Addr, c.Fd))
		}
	}

}

// SendMessage Message distribution
func (c *Client) SendMessage(message []byte) {

	if c == nil {
		logger.ErrorString("Websocket", "SendMessage",
			fmt.Sprintf("Not found websocket client, fd: %s", c.Fd))
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
	c.OwnerChannel.LoadOrStore(channel, pubSub)
	c.receive <- pubSub
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
func (c *Client) GetAllChan() (channels []string, pubSubs []*redis.PubSub) {
	c.Channel.Range(func(key, value any) bool {
		pubSub, ok := value.(*redis.PubSub)
		if ok {
			channels = append(channels, key.(string))
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
	c.receive <- pubSub
	return
}

// Receive subscription messages
func (c *Client) Receive() {
	for {
		select {
		case pubSub := <-c.receive:
			if pubSub != nil {
				go func(sub *redis.PubSub) {
					ch := sub.Channel()
					for message := range ch {
						c.Send <- []byte(message.Payload)
					}
				}(pubSub)
			}
		case <-c.close:
			return
		}
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

func (c *Client) CloseSubscribe(channel string) (err error) {
	value, ok := c.OwnerChannel.Load(channel)
	if !ok {
		return errors.New("not found channel")
	}
	pubSub, ok := value.(*redis.PubSub)
	if !ok {
		return errors.New("channel type error")
	}
	err = pubSub.Close()
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

// Heartbeat The scheduled task clears timeout links
func (c *Client) Heartbeat() {
	gap := get.Int64("app.heartbeat_check_time", 1000)
	EventListener(time.Millisecond*time.Duration(gap), func() {
		if c.IsHeartbeatTimeout(time.Now().Unix()) {
			c.Close()
		}
	}, c.close)
}
