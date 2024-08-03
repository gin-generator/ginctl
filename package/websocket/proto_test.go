package websocket

import (
	"fmt"
	"github.com/gin-generator/ginctl/package/websocket/pb"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"net/http"
	"testing"
)

func TestWebSocketConnection(t *testing.T) {
	// 连接到WebSocket服务器
	url := "ws://127.0.0.1:9503/ws"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("连接WebSocket服务器失败: %v", err)
	}
	defer conn.Close()

	// 创建protobuf消息
	msg := &pb.Message{
		Event:   "ping",
		Request: "test_request",
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		t.Fatalf("protobuf消息序列化失败: %v", err)
	}

	// 发送protobuf消息
	if err := conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		t.Fatalf("发送消息失败: %v", err)
	}

	// 读取响应
	_, response, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("读取消息失败: %v", err)
	}

	// 解析响应
	var res pb.Response
	if err := proto.Unmarshal(response, &res); err != nil {
		t.Fatalf("响应解析失败: %v", err)
	}

	// 验证响应
	if res.Code != http.StatusOK {
		t.Fatalf("响应代码不符合预期: %v", res.Code)
	}
	if res.Message != "success" {
		t.Fatalf("响应消息不符合预期: %v", res.Message)
	}
	fmt.Println("响应内容:", res.Content)
}
