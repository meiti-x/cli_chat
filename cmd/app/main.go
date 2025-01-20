package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/meiti-x/snapp_chal/config"
	"github.com/meiti-x/snapp_chal/internal/chatroom"
	nats2 "github.com/meiti-x/snapp_chal/internal/nats"
	redis2 "github.com/meiti-x/snapp_chal/internal/redis"
	"github.com/meiti-x/snapp_chal/pkg/app_errors"
	"github.com/meiti-x/snapp_chal/pkg/events"
	"github.com/meiti-x/snapp_chal/pkg/logger"
	"github.com/meiti-x/snapp_chal/socket"
	"github.com/nats-io/nats.go"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	ctx      = context.Background() // Context for Redis operations
)

const ChatroomNameQuery = "chatroom"

// TODO: dockerize project
// TODO: add documents
// TODO: add git hook
// TODO: change structure of project
// TODO: add more commands(my message)
// TODO: clear online users in redis on server stop
// TODO: rename project
// TODO add logger

func main() {
	configPath := flag.String("c", "config.yml", "Path to the configuration file")
	flag.Parse()

	conf, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalln(fmt.Errorf("load config error: %w", err))
	}

	lo := logger.NewAppLogger(conf)
	lo.InitLogger(conf.Logger.Path)

	lo.Error("لاگ کوچکنرززززز:", err)
	nc, err := nats2.MustInitNats(conf)
	if err != nil {
		panic(app_errors.ErrNatsInit)
	}
	defer nc.Close()

	rdb := redis2.MustInitRedis(conf)
	defer func(rdb *redis.Client) {
		err := rdb.Close()
		if err != nil {
			log.Println(app_errors.ErrRedisClose)
		}
	}(rdb)

	// TODO: add simple logger pkg
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(app_errors.ErrSocketUpgradeFailed, err)
			return
		}
		defer conn.Close()

		chatroomName := r.URL.Query().Get(ChatroomNameQuery)
		if chatroomName == "" {
			chatroomName = "general"
		}

		clientIP, _, _ := net.SplitHostPort(r.RemoteAddr)

		// Add user to Redis set for the specific chatroom
		onlineUsersKey := fmt.Sprintf("chatroom:%s:online_users", chatroomName)
		if err := rdb.SAdd(ctx, onlineUsersKey, clientIP).Err(); err != nil {
			log.Println(app_errors.ErrRedisOperationFailed, err)
			return
		}

		// Notify chat about the new user
		if err := chatroom.SendWelcomeMessage(clientIP, conn); err != nil {
			log.Printf(app_errors.ErrSendWelcomeMessage.Error(), err)
		}

		totalUsers, _ := rdb.SCard(ctx, onlineUsersKey).Result()
		joinMessage := map[string]interface{}{
			"event":      events.EventUserJoined,
			"chatroom":   chatroomName,
			"ip":         clientIP,
			"totalUsers": totalUsers,
		}
		if err := chatroom.SendJoinRoomMessage(joinMessage, nc, chatroomName); err != nil {
			log.Println(app_errors.ErrSendJoinMessage, err)
		}
		if err = socket.BroadcastRedisUsers(rdb, conn, onlineUsersKey, joinMessage); err != nil {
			log.Println(app_errors.ErrSendOnlineUsers, err)
		}

		go func() {
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					log.Println(app_errors.ErrSocketReadFailed, err)
					break
				}

				msg := string(message)

				if msg == "#users" {
					handleUsersCommandRedis(conn, rdb, onlineUsersKey)
					continue
				}

				var userMessage map[string]interface{}
				if err := json.Unmarshal(message, &userMessage); err != nil {
					log.Println(app_errors.ErrParseJSON, err)
					continue
				}

				userMessage["ip"] = clientIP
				userMessage["chatroom"] = chatroomName
				userMessageJSON, _ := json.Marshal(userMessage)
				err = nc.Publish(fmt.Sprintf("chatroom.%s", chatroomName), userMessageJSON)
				if err != nil {
					log.Println(err)
				}
			}

			// on User disconnected
			if err := rdb.SRem(ctx, onlineUsersKey, clientIP).Err(); err != nil {
				log.Println(app_errors.ErrRedisOperationFailed, err)
				return
			}

			totalUsers, _ := rdb.SCard(ctx, onlineUsersKey).Result()
			leaveMessage := map[string]interface{}{
				"event":      events.EventUserLeft,
				"chatroom":   chatroomName,
				"ip":         clientIP,
				"totalUsers": totalUsers,
			}
			if err = chatroom.SendLeaveRoomMessage(leaveMessage, nc, chatroomName); err != nil {
				log.Println(app_errors.ErrSendLeaveMessage, err)
			}
			if err = socket.BroadcastRedisUsers(rdb, conn, onlineUsersKey, leaveMessage); err != nil {
				log.Println(app_errors.ErrSendOnlineUsers, err)
			}
		}()

		sub, _ := nc.SubscribeSync(fmt.Sprintf("chatroom.%s", chatroomName))
		defer sub.Unsubscribe()

		for {
			msg, err := sub.NextMsg(1 * time.Second)
			if err != nil && !errors.Is(err, nats.ErrTimeout) {
				log.Println(app_errors.ErrNATSReceivedFailed, err)
				return
			}
			if msg != nil {
				if err := conn.WriteMessage(websocket.TextMessage, msg.Data); err != nil {
					log.Println(app_errors.ErrSocketWriteFailed, err)
					return
				}
			}
		}
	})

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", conf.Server.Port),
	}
	go func() {
		fmt.Printf("Server started at %s:%d\n", conf.Server.Host, conf.Server.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	fmt.Println("Server gracefully stopped.")
}

// handleUsersCommandRedis sends the list of online users from Redis
func handleUsersCommandRedis(conn *websocket.Conn, rdb *redis.Client, subj string) {
	users, err := rdb.SMembers(ctx, subj).Result()
	if err != nil {
		log.Println(app_errors.ErrRedisOperationFailed, err)
		return
	}

	response := map[string]interface{}{
		"event": events.EventUserList,
		"users": users,
	}
	responseJSON, _ := json.Marshal(response)

	if err := conn.WriteMessage(websocket.TextMessage, responseJSON); err != nil {
		log.Println(app_errors.ErrSendOnlineUsers, err)
	}
}
