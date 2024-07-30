package websocket

import (
	"encoding/json"
	"github.com/gin-generator/ginctl/package/validator"
	"net/http"
)

type JsonHandler struct{}

type Response struct {
	Code    uint32 `json:"code"`
	Message string `json:"message"`
	Content string `json:"content"`
}

type Message struct {
	Event   string `json:"event" validate:"required"`
	Request string `json:"request"`
}

func NewJsonHandler() *JsonHandler {
	return &JsonHandler{}
}

func (j *JsonHandler) Distribute(client *Client, message []byte) (err error) {
	var msg Message
	err = json.Unmarshal(message, &msg)
	if err != nil {
		return
	}

	err = validator.ValidateStructWithOutCtx(msg)
	if err != nil {
		return j.Do(client, &Response{
			Code:    http.StatusBadRequest,
			Message: "event not found",
		})
	}

	handler, err := GetHandler(msg.Event)
	if err != nil {
		return j.Do(client, &Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	response := handler(&Request{
		Client: client,
		Send:   []byte(msg.Request),
	})

	if response != nil {
		return j.Do(client, response)
	}
	return
}

func (j *JsonHandler) Do(client *Client, response *Response) (err error) {
	bytes, err := json.Marshal(response)
	if err != nil {
		return
	}

	client.SendMessage(bytes)
	return
}
