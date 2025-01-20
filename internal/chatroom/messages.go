package chatroom

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
)

func SendWelcomeMessage(clientIP string, conn *websocket.Conn) error {
	welcomeMessage := map[string]interface{}{
		"event":   "welcome",
		"message": fmt.Sprintf("Welcome to the chatroom! Your IP: %s", clientIP),
		"instructions": []string{
			"Type your message and hit Enter to chat.",
			"Use special commands like #help to see all avail",
			"Enjoy your stay and follow chatroom etiquette!",
		},
	}
	welcomeMessageJSON, _ := json.Marshal(welcomeMessage)
	if err := conn.WriteMessage(websocket.TextMessage, welcomeMessageJSON); err != nil {
		return err
	}
	return nil
}

func SendJoinRoomMessage(msg map[string]interface{}, nc *nats.Conn, subj string) error {
	joinMessageJSON, _ := json.Marshal(msg)
	if err := nc.Publish(fmt.Sprintf("chatroom.%s", subj), joinMessageJSON); err != nil {
		return err
	}
	return nil
}

func SendLeaveRoomMessage(msg map[string]interface{}, nc *nats.Conn, subj string) error {
	leaveMessageJSON, _ := json.Marshal(msg)
	if err := nc.Publish(fmt.Sprintf("chatroom.%s", subj), leaveMessageJSON); err != nil {
		return err
	}
	return nil
}
