package middleware

import (
	"github.com/gin-generator/ginctl/package/websocket"
)

func Auth(next websocket.DisposeFunc) websocket.DisposeFunc {
	return func(req *websocket.Request) *websocket.Response {
		// TODO 测试
		if req.Client.Fd != "" {
			return &websocket.Response{
				Code:    401,
				Message: "auth failed",
			}
		}
		return next(req)
	}
}
