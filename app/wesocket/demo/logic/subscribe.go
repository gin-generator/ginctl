package logic

import (
	"encoding/json"
	"github.com/gin-generator/ginctl/package/websocket"
)

type SubscribeRequest struct {
	Channel string
}

func Subscribe(request *websocket.Request) (response *websocket.Response) {
	response = &websocket.Response{}
	req := &SubscribeRequest{}
	err := json.Unmarshal(request.Send, req)
	if err != nil {
		response.Code = 500
		response.Message = err.Error()
	}

	err = request.Client.Subscribe(req.Channel)
	if err != nil {
		response.Code = 500
		response.Message = err.Error()
	}

	response.Code = 200
	response.Message = "success"
	return
}
