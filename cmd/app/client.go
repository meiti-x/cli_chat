package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/meiti-x/snapp_task/config"
	"github.com/meiti-x/snapp_task/internal/models"
	"github.com/meiti-x/snapp_task/pkg/adapters/storage"
	"github.com/meiti-x/snapp_task/pkg/app_errors"
	db2 "github.com/meiti-x/snapp_task/pkg/db"
	"github.com/meiti-x/snapp_task/pkg/logger"
	"github.com/meiti-x/snapp_task/usecase"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	serverURL := flag.String("server", "ws://localhost:8080/ws", "WebSocket server URL")
	chatroom := flag.String("chatroom", "general", "Chatroom to join")
	configPath := flag.String("c", "config.yml", "Path to the configuration file")
	flag.Parse()

	conf, err := config.LoadConfig(*configPath)

	logger := logger.NewAppLogger(conf)
	logger.InitLogger(conf.Logger.Path)

	db, err := db2.InitDB()
	if err != nil {
		logger.Error(app_errors.ErrInitDB)
	}
	userRepo := storage.NewUserRepository(db)
	messageRepo := storage.NewMessageRepository(db)

	// User authentication
	fmt.Println("Welcome! Please log in or register:")
	authenticatedUser := usecase.AuthenticateUser(userRepo)

	connURL := fmt.Sprintf("%s?chatroom=%s", *serverURL, *chatroom)

	conn, _, err := websocket.DefaultDialer.Dial(connURL, nil)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	fmt.Printf("Connected to chatroom '%s' as '%s'. Type messages or use the '#users' command to see online users.\n", *chatroom, authenticatedUser)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Channel to receive messages from the server
	messageChan := make(chan string)

	// Goroutine to listen
	go func() {
		defer close(messageChan)
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading message: %v", err)
				return
			}
			messageChan <- string(msg)
		}
	}()

	// Goroutine to send messages to the server
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text := scanner.Text()

			// Check for special commands
			if text == "#users" {
				if err := conn.WriteMessage(websocket.TextMessage, []byte("#users")); err != nil {
					log.Printf("Error sending command: %v", err)
				}
				continue
			}

			// Parse JSON input to extract "message" field
			var parsedMessage map[string]string
			if err := json.Unmarshal([]byte(text), &parsedMessage); err != nil {
				log.Printf("Invalid message format: %v", err)
				continue
			}

			messageContent, ok := parsedMessage["message"]
			if !ok {
				log.Println("Message field not found in input")
				continue
			}

			chatMessage := &models.Message{
				Username:  authenticatedUser,
				Chatroom:  *chatroom,
				Content:   messageContent,
				CreatedAt: time.Now(),
			}
			if err := messageRepo.CreateMessage(chatMessage); err != nil {
				log.Printf("Failed to save message to database: %v", err)
				continue
			}

			// Send the original JSON message to the server
			if err := conn.WriteMessage(websocket.TextMessage, []byte(text)); err != nil {
				log.Printf("Error sending message: %v", err)
			}
		}
	}()

	// Main loop to handle signals
	for {
		select {
		case <-stop:
			fmt.Println("\nDisconnecting...")
			return
		case msg := <-messageChan:
			if msg == "" || len(msg) == 0 {
				continue
			}
			fmt.Printf("Server: %s\n", msg)
		}
	}
}
