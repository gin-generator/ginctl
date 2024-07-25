package logic

import (
	"encoding/json"
	"github.com/gin-generator/ginctl/package/websocket"
)

type BroadcastRequest struct {
	Channel string
	Message string
}

func Broadcast(request *websocket.Request) (response *websocket.Response) {
	response = &websocket.Response{}
	req := &BroadcastRequest{}

	err := json.Unmarshal(request.Send, req)
	if err != nil {
		response.Code = 500
		response.Message = err.Error()
	}

	err = request.Client.Publish(req.Channel, []byte(req.Message))
	if err != nil {
		response.Code = 500
		response.Message = err.Error()
	}
	response.Code = 200
	response.Message = "success"
	return
}
