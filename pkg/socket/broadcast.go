package socket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	c "github.com/meiti-x/cli_chat/pkg/cache"
	"log"
)

// BroadcastRedisUsers sends a message to all online users in the chatroom
func BroadcastRedisUsers(rdb c.Provider, ws *websocket.Conn, onlineUsersKey string, message map[string]interface{}) error {
	ctx := context.Background()

	onlineUsers, err := rdb.GetSetMembers(ctx, onlineUsersKey)
	fmt.Println(onlineUsers)
	if err != nil {
		return err
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}

	for _, user := range onlineUsers {

		if ws == nil {
			log.Printf("No WebSocket connection found for user: %s\n", user)
			continue
		}

		err := ws.WriteMessage(websocket.TextMessage, messageJSON)
		if err != nil {
			fmt.Println(err)
			log.Printf("Failed to send message to user: %s, error: %v\n", user, err)
			continue
		}
	}

	return nil
}
