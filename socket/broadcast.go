package socket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"log"
)

func BroadcastRedisUsers(rdb *redis.Client, ws *websocket.Conn, onlineUsersKey string, message map[string]interface{}) error {
	ctx := context.Background()

	onlineUsers, err := rdb.SMembers(ctx, onlineUsersKey).Result()
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
