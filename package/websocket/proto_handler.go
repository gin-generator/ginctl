package websocket

import (
	"github.com/gin-generator/ginctl/package/websocket/pb"
	"google.golang.org/protobuf/proto"
	"net/http"
)

type ProtoHandler struct{}

func NewProtoHandler() *ProtoHandler {
	return &ProtoHandler{}
}

func (p *ProtoHandler) Distribute(client *Client, message []byte) (err error) {

	var msg pb.Message
	err = proto.Unmarshal(message, &msg)
	if err != nil {
		return
	}

	if msg.Event == "" {
		return p.Do(client, &pb.Response{
			Code:    http.StatusBadRequest,
			Message: "event not found",
		})
	}

	handler, err := GetHandler(msg.Event)
	if err != nil {
		return p.Do(client, &pb.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	response := handler(&Request{
		Client: client,
		Send:   []byte(msg.Request),
	})

	if response != nil {
		return p.Do(client, &pb.Response{
			Code:    response.Code,
			Message: response.Message,
			Content: response.Content,
		})
	}
	return
}

func (p *ProtoHandler) Do(client *Client, response *pb.Response) (err error) {
	bytes, err := proto.Marshal(response)
	if err != nil {
		return
	}

	client.SendMessage(bytes)
	return
}
