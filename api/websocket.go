//package api
//
//import (
//	"encoding/json"
//	"fmt"
//	"github.com/gorilla/websocket"
//	"github.com/nats-io/nats.go"
//	"log"
//	"net"
//	"net/http"
//	"sync"
//	"time"
//)
//
//var (
//	upgrader    = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
//	onlineUsers = make(map[*websocket.Conn]string) // Track online users
//	mu          sync.Mutex                         // Protect onlineUsers map
//)
//
//func RegisterSocket() {
//	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
//		conn, err := upgrader.Upgrade(w, r, nil)
//		if err != nil {
//			log.Printf("WebSocket upgrade error: %v", err)
//			return
//		}
//		defer conn.Close()
//
//		// Get client IP
//		clientIP, _, _ := net.SplitHostPort(r.RemoteAddr)
//
//		// Add user to the online user list
//		mu.Lock()
//		onlineUsers[conn] = clientIP
//		totalUsers := len(onlineUsers)
//		mu.Unlock()
//
//		// Notify about the new user
//
//		// Send a welcome message to the joining user
//		welcomeMessage := map[string]interface{}{
//			"event":   "welcome",
//			"message": fmt.Sprintf("Welcome to the chatroom! Your IP: %s", clientIP),
//			"instructions": []string{
//				"Type your message and hit Enter to chat.",
//				"Use special commands like #help to see all avail",
//				"Enjoy your stay and follow chatroom etiquette!",
//			},
//		}
//		welcomeMessageJSON, _ := json.Marshal(welcomeMessage)
//		if err := conn.WriteMessage(websocket.TextMessage, welcomeMessageJSON); err != nil {
//			log.Printf("Error sending welcome message: %v", err)
//		}
//
//		joinMessage := map[string]interface{}{
//			"event":      "user_joined",
//			"ip":         clientIP,
//			"totalUsers": totalUsers,
//		}
//		joinMessageJSON, _ := json.Marshal(joinMessage)
//		nc.Publish("chatroom", joinMessageJSON) // Publish join event to NATS
//
//		// Broadcast join message to all users
//		broadcast(nc, joinMessage)
//
//		// Handle incoming WebSocket messages
//		go func() {
//			for {
//				_, message, err := conn.ReadMessage() // Read raw WebSocket message
//				if err != nil {
//					log.Printf("Error reading WebSocket message: %v", err)
//					break
//				}
//
//				msg := string(message) // Convert the raw message to a string
//
//				// Check for special commands
//				if msg == "#users" {
//					handleUsersCommand(conn)
//					continue
//				}
//
//				// Try to parse the message as JSON for regular messages
//				var userMessage map[string]interface{}
//				if err := json.Unmarshal(message, &userMessage); err != nil {
//					log.Printf("Error parsing JSON message: %v", err)
//					continue
//				}
//
//				// Add IP to the message and publish to NATS
//				userMessage["ip"] = clientIP
//				userMessageJSON, _ := json.Marshal(userMessage)
//				nc.Publish("chatroom", userMessageJSON)
//			}
//
//			// User disconnected
//			mu.Lock()
//			delete(onlineUsers, conn)
//			totalUsers := len(onlineUsers)
//			mu.Unlock()
//
//			// Notify about the user leaving
//			leaveMessage := map[string]interface{}{
//				"event":      "user_left",
//				"ip":         clientIP,
//				"totalUsers": totalUsers,
//			}
//			leaveMessageJSON, _ := json.Marshal(leaveMessage)
//			nc.Publish("chatroom", leaveMessageJSON) // Publish leave event to NATS
//
//			// Broadcast leave message to all users
//			broadcast(nc, leaveMessage)
//		}()
//
//		// Forward NATS messages to WebSocket
//		sub, _ := nc.SubscribeSync("chatroom")
//		defer sub.Unsubscribe()
//
//		for {
//			msg, err := sub.NextMsg(1 * time.Second)
//			if err != nil && err != nats.ErrTimeout {
//				log.Printf("Error receiving NATS message: %v", err)
//				return
//			}
//			if msg != nil {
//				if err := conn.WriteMessage(websocket.TextMessage, msg.Data); err != nil {
//					log.Printf("Error writing WebSocket message: %v", err)
//					return
//				}
//			}
//		}
//	})
//
//}
