package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type subscriptionMessage struct {
	Type    string          `json:"type"`
	Id      string          `json:"id,omitempty"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

func ConnectToGraphQLServer(url string) (*websocket.Conn, error) {
	conn, _, err := websocket.Dial(context.Background(), url, nil)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func Subscribe(conn *websocket.Conn, subscription string) error {
	msg := subscriptionMessage{
		Type: "connection_init",
	}

	err := wsjson.Write(context.Background(), conn, msg)
	if err != nil {
		return err
	}

	msg = subscriptionMessage{
		Type: "start",
		Id:   "1",
		Payload: json.RawMessage(fmt.Sprintf(
			`{"query": %q}`,
			subscription,
		)),
	}

	err = wsjson.Write(context.Background(), conn, msg)
	if err != nil {
		return err
	}

	return nil
}

func ReceiveMessages(conn *websocket.Conn, next func(payload json.RawMessage)) {
	for {
		var msg subscriptionMessage
		err := wsjson.Read(context.Background(), conn, &msg)
		if err != nil {
			log.Println("Error receiving message:", err)
			return
		}

		if msg.Type == "data" {
			next(msg.Payload)
			//fmt.Println("Received data:", string(msg.Payload))
		}
		/*else {
			fmt.Println("Received message:", msg)
		}*/
	}
}
