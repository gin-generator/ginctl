package logic

import "github.com/gin-generator/ginctl/package/websocket"

func Ping(_ *websocket.Request) (response *websocket.Response) {
	response = websocket.NewResponse()
	response.Code = 200
	response.Message = "success"
	response.Content = "pong"
	return
}
