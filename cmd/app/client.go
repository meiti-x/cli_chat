package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/meiti-x/snapp_task/internal/models"
	db2 "github.com/meiti-x/snapp_task/pkg/db"
	"gorm.io/gorm"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	serverURL := flag.String("server", "ws://localhost:8080/ws", "WebSocket server URL")
	chatroom := flag.String("chatroom", "general", "Chatroom to join")
	flag.Parse()

	db, err := db2.InitDB()

	// User authentication
	authenticatedUser := authenticateUser(db)

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
			if err := db.Create(&chatMessage).Error; err != nil {
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
			fmt.Printf("Server: %s\n", msg)
		}
	}
}

// authenticateUser handles user login and registration.
func authenticateUser(db *gorm.DB) string {
	fmt.Println("Welcome! Please log in or register:")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("1. Login")
		fmt.Println("2. Register")
		fmt.Print("Choose an option: ")

		choice, _ := reader.ReadString('\n')
		choice = choice[:len(choice)-1]

		if choice == "1" {
			return loginUser(db)
		} else if choice == "2" {
			return registerUser(db)
		} else {
			fmt.Println("Invalid option. Please try again.")
		}
	}
}

// loginUser handles user login.
func loginUser(db *gorm.DB) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter username: ")
	username, _ := reader.ReadString('\n')
	username = username[:len(username)-1]

	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')
	password = password[:len(password)-1]

	var user *models.User
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		fmt.Println("User not found. Please register.")
		return authenticateUser(db)
	}

	// Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		fmt.Println("Invalid password. Please try again.")
		return authenticateUser(db)
	}

	fmt.Println("Login successful!")
	return username
}

// registerUser handles new user registration.
func registerUser(db *gorm.DB) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter username: ")
	username, _ := reader.ReadString('\n')
	username = username[:len(username)-1]

	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')
	password = password[:len(password)-1]

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// Save the user to the database
	user := models.User{Username: username, Password: string(hashedPassword)}
	if err := db.Create(&user).Error; err != nil {
		fmt.Printf("Failed to register user: %v\n", err)
		return authenticateUser(db)
	}

	fmt.Println("Registration successful!")
	return username
}
