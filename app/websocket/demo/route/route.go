package route

import (
	"github.com/gin-generator/ginctl/app/websocket/demo/logic"
	m "github.com/gin-generator/ginctl/app/websocket/demo/middleware"
	"github.com/gin-generator/ginctl/package/websocket"
)

func RegisterDemoRoute() {
	r := websocket.NewRouter()
	// websocket event
	r.Register("ping", logic.Ping, m.Auth)
	r.Register("heartbeat", logic.Heartbeat)
	r.Register("create_channel", logic.CreateChannel, m.Auth)
	r.Register("subscribe", logic.Subscribe, m.Auth)
	r.Register("broadcast", logic.Broadcast)
	r.Register("staff_broadcast", logic.StaffBroadcast)
}
