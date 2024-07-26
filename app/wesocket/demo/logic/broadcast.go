package logic

import (
	"encoding/json"
	"github.com/gin-generator/ginctl/package/validator"
	"github.com/gin-generator/ginctl/package/websocket"
)

type BroadcastRequest struct {
	Channel string `validate:"required"`
	Message string `validate:"required"`
}

func Broadcast(request *websocket.Request) (response *websocket.Response) {
	response = &websocket.Response{}
	req := &BroadcastRequest{}

	err := json.Unmarshal(request.Send, req)
	if err != nil {
		response.Code = 500
		response.Message = err.Error()
		return
	}

	err = validator.ValidateStructWithOutCtx(req)
	if err != nil {
		response.Code = 500
		response.Message = err.Error()
		return
	}

	err = request.Client.Publish(req.Channel, []byte(req.Message))
	if err != nil {
		response.Code = 500
		response.Message = err.Error()
		return
	}
	response.Code = 200
	response.Message = "success"
	return
}
