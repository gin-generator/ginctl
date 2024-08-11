package main

import (
	"fmt"
	"github.com/gin-generator/ginctl/app/websocket/demo/route"
	"github.com/gin-generator/ginctl/bootstrap"
	"github.com/gin-generator/ginctl/package/get"
	"github.com/gin-generator/ginctl/package/websocket"
	"net/http"
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
	websocket.NewClientManager()

	// Start the websocket serve
	http.HandleFunc("/ws", websocket.Upgrade)
	host := get.String("app.host", "127.0.0.1")
	port := get.Uint("app.port", 9503)
	fmt.Println(fmt.Sprintf("websocket demo serve start: %s:%d...", host, port))
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil)
	if err != nil {
		fmt.Println("websocket demo start failure")
	}
}
