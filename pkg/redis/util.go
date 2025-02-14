package redis

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/meiti-x/cli_chat/pkg/app_errors"
	c "github.com/meiti-x/cli_chat/pkg/cache"
	"github.com/meiti-x/cli_chat/pkg/events"
	"log"
)

// HandleUsersCommandRedis sends the list of online users from Redis
func HandleUsersCommandRedis(ctx context.Context, conn *websocket.Conn, rdb c.Provider, subj string) {
	users, err := rdb.GetSetMembers(ctx, subj)
	if err != nil {
		log.Println(err)
		return
	}

	// TODO: create a response pkg
	response := map[string]interface{}{
		"event": events.EventUserList,
		"users": users,
	}
	responseJSON, _ := json.Marshal(response)

	if err := conn.WriteMessage(websocket.TextMessage, responseJSON); err != nil {
		log.Println(app_errors.ErrSendOnlineUsers, err)
	}
}
