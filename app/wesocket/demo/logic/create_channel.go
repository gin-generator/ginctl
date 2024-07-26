package logic

import "github.com/gin-generator/ginctl/package/websocket"

func CreateChannel(request *websocket.Request) (response *websocket.Response) {
	response = &websocket.Response{}
	channel, err := request.Client.CreatChan()
	if err != nil {
		response.Code = 500
		response.Message = err.Error()
		return
	}
	response.Code = 200
	response.Message = "success"
	response.Data = channel
	return
}
