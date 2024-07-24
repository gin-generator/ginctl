package websocket

import (
	"fmt"
	"github.com/gin-generator/ginctl/package/get"
	"github.com/gin-generator/ginctl/package/logger"
	t "github.com/gin-generator/ginctl/package/time"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Client struct {
	Fd            string          // 每个连接唯一标识
	Addr          string          // 客户端ip地址
	Socket        *websocket.Conn // 用户连接
	Send          chan []byte     // 待发送的数据
	FirstTime     int64           // 首次连接事件
	HeartbeatTime int64           // 用户上次心跳时间
	Timeout       int64           // 自动断连时间
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

func (c *Client) SetHeartbeatTime(currentTime int64) {
	c.HeartbeatTime = currentTime
}

func (c *Client) IsHeartbeatTimeout(currentTime int64) (timeout bool) {
	if c.HeartbeatTime+c.Timeout <= currentTime {
		timeout = true
	}
	return
}
