package main

import (
	"fmt"
	"github.com/gin-generator/ginctl/app/wesocket/demo/route"
	"github.com/gin-generator/ginctl/bootstrap"
	"github.com/gin-generator/ginctl/package/get"
	ws "github.com/gin-generator/ginctl/package/websocket"
	"github.com/gorilla/websocket"
	"net/http"
)

var (
	manager *ws.ClientManager
)

func main() {
	// Load config.
	get.NewViper("env.yaml", "./etc")

	// Start basic server.
	bootstrap.SetupLogger()
	bootstrap.SetupDB()
	bootstrap.SetupRedis()

	// Register event
	route.RegisterDemoRoute()

	// Start the websocket scheduler
	manager = ws.NewClientManager()
	go manager.Scheduler()

	// Start the websocket serve
	http.HandleFunc("/ws", upgrade)
	host := get.String("app.host", "127.0.0.1")
	port := get.Uint("app.port", 9503)
	fmt.Println(fmt.Sprintf("websocket demo serve start: %s:%d...", host, port))
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil)
	if err != nil {
		fmt.Println("websocket demo start failure")
	}
}

func upgrade(w http.ResponseWriter, req *http.Request) {

	if manager.Total+1 > manager.Max {
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
		fmt.Println("websocket demo start failure:", err.Error())
		http.NotFound(w, req)
		return
	}

	client := ws.NewClient(conn.RemoteAddr().String(), conn)
	client.Timeout = get.Int64("app.heartbeat_timeout", 600)
	client.Send <- []byte(fmt.Sprintf("{\"code\": 200,\"message\": \"success\",\"data\": {\"fd\": \"%s\"}}", client.Fd))

	// 监听读
	go client.Read()
	// 监听写
	go client.Write()
	// 监听订阅
	go client.Receive()
	// 心跳检测
	go client.Heartbeat()

	manager.Register <- client
}
