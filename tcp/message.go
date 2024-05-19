package tcp

import (
	"encoding/json"
	"errors"
	"time"
)

const (
	maxMessageLength = 2048 // max lengths
	readTimeout      = time.Second * 3
	writeTimeout     = time.Second * 3
)

const (
	RESPONSE = 0
	COMMAND  = 1
	FILE     = 2
)

const (
	RESPONSE_DESCRIPTION = "RESPONSE_DESCRIPTION"
	RESPONSE_TABLE       = "RESPONSE_TABLE"
)

type Message struct {
	Type    uint
	Command []string
	Data    []byte
}

func NewMessage(typ uint, cmd []string, data []byte) *Message {
	return &Message{
		Type:    typ,
		Command: cmd,
		Data:    data,
	}
}

func NewResponse(data map[string]interface{}) *Message {
	msg := &Message{
		Type:    RESPONSE,
		Command: nil,
		Data:    nil,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil
	}
	byteData := []byte(jsonData)
	msg.Data = byteData
	return msg
}

func NewResponseText(text string) *Message {
	msg := &Message{
		Type:    RESPONSE,
		Command: nil,
		Data:    nil,
	}
	data := map[string]interface{}{
		"status":      "OK",
		"description": text,
		"type":        RESPONSE_DESCRIPTION,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil
	}
	byteData := []byte(jsonData)
	msg.Data = byteData
	return msg

}

func (msg *Message) GetResponse() (map[string]interface{}, error) {

	var result map[string]interface{}

	err := json.Unmarshal(msg.Data, &result)
	if err != nil {
		return nil, err
	}

	return result, nil

}

func (msg *Message) GetCommand() string {
	return msg.Command[0]
}

func (msg *Message) GetCommandArgs(pos int) (string, error) {
	if len(msg.Command) >= pos {
		return msg.Command[pos], nil
	} else {
		return "", errors.New("Parameter not found")
	}

}
