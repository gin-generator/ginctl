package logic

import (
	"github.com/gin-generator/ginctl/package/websocket"
	"time"
)

func Heartbeat(request *websocket.Request) (response *websocket.Response) {
	response = &websocket.Response{}
	request.Client.HeartbeatTime = time.Now().Unix()
	response.Code = 200
	response.Message = "success"
	response.Data = "pong"
	return
}
