package logic

import "github.com/gin-generator/ginctl/package/websocket"

func Ping(_ *websocket.Request) (response *websocket.Response) {
	response = &websocket.Response{}
	response.Code = 200
	response.Message = "success"
	response.Data = "pong"
	return
}
