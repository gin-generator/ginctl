package logic

import (
	"encoding/json"
	"github.com/gin-generator/ginctl/package/validator"
	"github.com/gin-generator/ginctl/package/websocket"
)

type SubscribeRequest struct {
	Channel string `json:"channel" validate:"required"`
}

func Subscribe(request *websocket.Request) (response *websocket.Response) {
	response = &websocket.Response{}
	req := &SubscribeRequest{}
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

	err = request.Client.Subscribe(req.Channel)
	if err != nil {
		response.Code = 500
		response.Message = err.Error()
		return
	}

	response.Code = 200
	response.Message = "success"
	return
}
