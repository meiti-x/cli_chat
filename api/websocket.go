package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/meiti-x/snapp_task/internal/chatroom"
	"github.com/meiti-x/snapp_task/pkg/app_errors"
	c "github.com/meiti-x/snapp_task/pkg/cache"
	"github.com/meiti-x/snapp_task/pkg/events"
	"github.com/meiti-x/snapp_task/pkg/redis"
	"github.com/meiti-x/snapp_task/pkg/socket"
	"github.com/nats-io/nats.go"
	"net"
	"net/http"
	"time"
)

const ChatroomNameQuery = "chatroom"

var (
	upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	ctx      = context.Background() // Context for Redis operations
)

func InitWS(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			s.logger.Error(app_errors.ErrSocketUpgradeFailed, err)
			return
		}
		defer conn.Close()
		defer s.rdb.CloseConnection()
		defer s.nats.Close()

		chatroomName := r.URL.Query().Get(ChatroomNameQuery)
		if chatroomName == "" {
			chatroomName = "general"
		}

		clientIP, _, _ := net.SplitHostPort(r.RemoteAddr)

		// Add user to Redis set for the specific chatroom
		onlineUsersKey := fmt.Sprintf("chatroom:%s:online_users", chatroomName)
		if err := s.rdb.AddSetMember(ctx, onlineUsersKey, clientIP); err != nil {
			s.logger.Error(err)
			return
		}
		s.logger.Info(clientIP, " is joined ", chatroomName)

		// Notify chat about the new user
		if err := chatroom.SendWelcomeMessage(clientIP, conn); err != nil {
			s.logger.Error(app_errors.ErrSendWelcomeMessage)
			return
		}

		totalUsers, _ := s.rdb.GetSetSize(ctx, onlineUsersKey)
		joinMessage := map[string]interface{}{
			"event":      events.EventUserJoined,
			"chatroom":   chatroomName,
			"ip":         clientIP,
			"totalUsers": totalUsers,
		}
		if err := chatroom.SendJoinRoomMessage(joinMessage, s.nats, chatroomName); err != nil {
			s.logger.Error(app_errors.ErrSendJoinMessage, err)
		}
		if err = socket.BroadcastRedisUsers(s.rdb, conn, onlineUsersKey, joinMessage); err != nil {
			s.logger.Error(app_errors.ErrSendOnlineUsers, err)
		}

		go func() {
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					s.logger.Error(app_errors.ErrSocketReadFailed, err)
					break
				}

				msg := string(message)
				if msg == "#users" {
					redis.HandleUsersCommandRedis(ctx, conn, s.rdb, onlineUsersKey)
					continue
				}

				var userMessage map[string]interface{}
				if err := json.Unmarshal(message, &userMessage); err != nil {
					s.logger.Error(app_errors.ErrParseJSON, err)
					continue
				}

				userMessage["ip"] = clientIP
				userMessage["chatroom"] = chatroomName
				userMessageJSON, _ := json.Marshal(userMessage)
				err = s.nats.Publish(fmt.Sprintf("chatroom.%s", chatroomName), userMessageJSON)
				if err != nil {
					s.logger.Error(err)
				}
			}

			// on User disconnected
			err := s.rdb.RemoveSetMember(ctx, onlineUsersKey, clientIP)
			if err != nil {
				s.logger.Error(c.ErrRedisOperationFailed)
				return
			}

			totalUsers, err := s.rdb.GetSetSize(ctx, onlineUsersKey)
			if err != nil {
				s.logger.Error(c.ErrRedisOperationFailed)
				return
			}

			leaveMessage := map[string]interface{}{
				"event":      events.EventUserLeft,
				"chatroom":   chatroomName,
				"ip":         clientIP,
				"totalUsers": totalUsers,
			}
			if err = chatroom.SendLeaveRoomMessage(leaveMessage, s.nats, chatroomName); err != nil {
				s.logger.Error(app_errors.ErrSendLeaveMessage, err)
			}
			if err = socket.BroadcastRedisUsers(s.rdb, conn, onlineUsersKey, leaveMessage); err != nil {
				s.logger.Error(app_errors.ErrSendOnlineUsers, err)
			}
		}()

		sub, _ := s.nats.SubscribeSync(fmt.Sprintf("chatroom.%s", chatroomName))
		defer sub.Unsubscribe()

		for {
			msg, err := sub.NextMsg(1 * time.Second)
			if err != nil && !errors.Is(err, nats.ErrTimeout) {
				s.logger.Error(app_errors.ErrNATSReceivedFailed, err)
				return
			}
			if msg != nil {
				if err := conn.WriteMessage(websocket.TextMessage, msg.Data); err != nil {
					s.logger.Error(app_errors.ErrSocketWriteFailed, err)
					return
				}
			}
		}
	}

}
