package websocket

import (
	"fmt"
	"github.com/gin-generator/ginctl/package/get"
	"github.com/gorilla/websocket"
	"net/http"
)

func Upgrade(w http.ResponseWriter, req *http.Request) {

	if Manager.Total+1 > Manager.Max {
		http.Error(w, "websocket service connections exceeded the upper limit", http.StatusInternalServerError)
		return
	}

	conn, err := (&websocket.Upgrader{
		ReadBufferSize:  get.Int("app.read_buffer_size", 4096),
		WriteBufferSize: get.Int("app.write_buffer_size", 4096),
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(w, req, nil)

	if err != nil {
		http.NotFound(w, req)
		return
	}

	client := NewClient(conn.RemoteAddr().String(), conn, req)
	client.Send <- []byte(fmt.Sprintf("{\"code\": 200,\"message\": \"success\",\"content\": \"%s\"}", client.Fd))

	// 监听读
	go client.Read()
	// 监听写
	go client.Write()
	// 监听订阅
	go client.Receive()
	// 心跳检测
	go client.Heartbeat()

	Manager.Register <- client
}
