package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/meiti-x/cli_chat/internal/chatroom"
	"github.com/meiti-x/cli_chat/internal/models"
	"github.com/meiti-x/cli_chat/pkg/adapters/storage"
	"github.com/meiti-x/cli_chat/pkg/app_errors"
	c "github.com/meiti-x/cli_chat/pkg/cache"
	"github.com/meiti-x/cli_chat/pkg/events"
	"github.com/meiti-x/cli_chat/pkg/redis"
	"github.com/meiti-x/cli_chat/pkg/socket"
	"github.com/nats-io/nats.go"
	"golang.org/x/crypto/bcrypt"
	"net"
	"net/http"
	"time"
)

const ChatroomNameQuery = "chatroom"

var (
	upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	ctx      = context.Background() // Context for Redis operations
)

// TODO: to much coupling here, its better expose to the usecase layer
func InitWS(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			s.Logger.Error(app_errors.ErrSocketUpgradeFailed, err)
			return
		}

		defer conn.Close()

		chatroomName := r.URL.Query().Get(ChatroomNameQuery)
		messageRepo := storage.NewMessageRepository(s.Db)

		if chatroomName == "" {
			chatroomName = "general"
		}

		clientIP, _, _ := net.SplitHostPort(r.RemoteAddr)

		// Add user to Redis set for the specific chatroom
		onlineUsersKey := fmt.Sprintf("chatroom:%s:online_users", chatroomName)
		if err := s.Rdb.AddSetMember(ctx, onlineUsersKey, clientIP); err != nil {
			s.Logger.Error(err)
			return
		}
		s.Logger.Info(clientIP, " is joined ", chatroomName)

		// Notify chat about the new user
		if err := chatroom.SendWelcomeMessage(clientIP, conn); err != nil {
			s.Logger.Error(app_errors.ErrSendWelcomeMessage)
			return
		}

		totalUsers, _ := s.Rdb.GetSetSize(ctx, onlineUsersKey)
		joinMessage := map[string]interface{}{
			"event":      events.EventUserJoined,
			"chatroom":   chatroomName,
			"ip":         clientIP,
			"totalUsers": totalUsers,
		}
		if err := chatroom.SendJoinRoomMessage(joinMessage, s.Nats, chatroomName); err != nil {
			s.Logger.Error(app_errors.ErrSendJoinMessage, err)
		}
		if err = socket.BroadcastRedisUsers(s.Rdb, conn, onlineUsersKey, joinMessage); err != nil {
			s.Logger.Error(app_errors.ErrSendOnlineUsers, err)
		}

		go func() {
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					s.Logger.Error(app_errors.ErrSocketReadFailed, err)
					break
				}

				msg := string(message)
				if msg == "#users" {
					redis.HandleUsersCommandRedis(ctx, conn, s.Rdb, onlineUsersKey)
					continue
				}

				var userMessage map[string]interface{}
				if err := json.Unmarshal(message, &userMessage); err != nil {
					s.Logger.Error(app_errors.ErrParseJSON, err)
					continue
				}

				userMessage["ip"] = clientIP
				userMessage["chatroom"] = chatroomName
				userMessageJSON, _ := json.Marshal(userMessage)
				err = s.Nats.Publish(fmt.Sprintf("chatroom.%s", chatroomName), userMessageJSON)
				if err != nil {
					s.Logger.Error(err)
				}

				err = messageRepo.CreateMessage(&models.Message{
					Username: "",
					Chatroom: chatroomName,
					Content:  msg,
				})

				if err != nil {
					s.Logger.Error(err)
				}

			}

			// on User disconnected
			err := s.Rdb.RemoveSetMember(ctx, onlineUsersKey, clientIP)
			if err != nil {
				s.Logger.Error(c.ErrRedisOperationFailed)
				return
			}

			totalUsers, err := s.Rdb.GetSetSize(ctx, onlineUsersKey)
			if err != nil {
				s.Logger.Error(c.ErrRedisOperationFailed)
				return
			}

			leaveMessage := map[string]interface{}{
				"event":      events.EventUserLeft,
				"chatroom":   chatroomName,
				"ip":         clientIP,
				"totalUsers": totalUsers,
			}
			if err = chatroom.SendLeaveRoomMessage(leaveMessage, s.Nats, chatroomName); err != nil {
				s.Logger.Error(app_errors.ErrSendLeaveMessage, err)
			}
			if err = socket.BroadcastRedisUsers(s.Rdb, conn, onlineUsersKey, leaveMessage); err != nil {
				s.Logger.Error(app_errors.ErrSendOnlineUsers, err)
			}
		}()

		sub, _ := s.Nats.SubscribeSync(fmt.Sprintf("chatroom.%s", chatroomName))
		defer sub.Unsubscribe()

		for {
			msg, err := sub.NextMsg(1 * time.Second)
			if err != nil && !errors.Is(err, nats.ErrTimeout) {
				s.Logger.Error(app_errors.ErrNATSReceivedFailed, err)
				return
			}
			if msg != nil {
				if err := conn.WriteMessage(websocket.TextMessage, msg.Data); err != nil {
					s.Logger.Error(app_errors.ErrSocketWriteFailed, err)
					return
				}
			}
		}
	}

}
func LoginHandler(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var loginRequest struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		username := models.Username(loginRequest.Username)
		password := loginRequest.Password

		// Fetch user by username
		userRepo := storage.NewUserRepository(s.Db)
		user, err := userRepo.GetUserByUsername(&username)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		// FIXME: create util pkg for hash and compare password
		// Compare hashed password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Login successful",
		})
	}
}
func RegisterHandler(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if user.Username == "" || user.Password == "" {
			http.Error(w, "Username and password are required", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}
		user.Password = string(hashedPassword)

		userRepo := storage.NewUserRepository(s.Db)
		if err := userRepo.CreateUser(&user); err != nil {
			http.Error(w, "Failed to create user, username might already exist", http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "User registered successfully",
		})
	}
}
