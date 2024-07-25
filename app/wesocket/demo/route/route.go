package route

import (
	"github.com/gin-generator/ginctl/app/wesocket/demo/logic"
	"github.com/gin-generator/ginctl/package/websocket"
)

func RegisterDemoRoute() {
	// websocket event
	websocket.Register("ping", logic.Ping)
	websocket.Register("heartbeat", logic.Heartbeat)
	websocket.Register("create_channel", logic.CreateChannel)
	websocket.Register("subscribe", logic.Subscribe)
	websocket.Register("broadcast", logic.Broadcast)
}
