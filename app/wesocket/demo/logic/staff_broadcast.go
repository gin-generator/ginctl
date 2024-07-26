package logic

import (
	"github.com/gin-generator/ginctl/package/websocket"
)

func StaffBroadcast(request *websocket.Request) (response *websocket.Response) {
	response = &websocket.Response{}
	websocket.Manager.Broadcast <- []byte("注意注意！这里是全员广播！")
	response.Code = 200
	response.Message = "success"
	return
}
