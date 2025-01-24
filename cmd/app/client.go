package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/meiti-x/snapp_task/config"
	"github.com/meiti-x/snapp_task/pkg/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	// TODO: remove harcoded values
	serverURL := flag.String("server", "ws://localhost:8080/ws", "WebSocket server URL")
	chatroom := flag.String("chatroom", "general", "Chatroom to join")
	authURL := flag.String("auth", "http://localhost:8080/auth", "Authentication server URL")
	configPath := flag.String("c", "config.yml", "Path to the configuration file")
	flag.Parse()

	conf, err := config.LoadConfig(*configPath)

	logger := logger.NewAppLogger(conf)
	logger.InitLogger(conf.Logger.Path)

	// Prompt for login or registration
	if !authenticateUser(*authURL) {
		log.Fatal("Authentication failed. Exiting.")
		return
	}

	connURL := fmt.Sprintf("%s?chatroom=%s", *serverURL, *chatroom)

	conn, _, err := websocket.DefaultDialer.Dial(connURL, nil)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	fmt.Printf("Connected to chatroom '%s'. Type messages or use the '#users' command to see online users.\n", *chatroom)

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

// authenticateUser handles user login or registration.
func authenticateUser(authURL string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("1. Login")
		fmt.Println("2. Register")
		fmt.Print("Choose an option: ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		if choice == "1" {
			return loginUser(authURL)
		} else if choice == "2" {
			return registerUser(authURL)
		} else {
			fmt.Println("Invalid option. Please try again.")
		}
	}
}

// loginUser prompts the user for login credentials and authenticates with the server.
func loginUser(authURL string) bool {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	payload := map[string]string{
		"username": username,
		"password": password,
	}
	data, _ := json.Marshal(payload)

	resp, err := http.Post(fmt.Sprintf("%s/login", authURL), "application/json", strings.NewReader(string(data)))
	if err != nil {
		fmt.Println("Login failed. Please try again.")
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Login successful!")
		return true
	}
	fmt.Printf("Login failed with status: %d\n", resp.StatusCode)
	return false
}

// registerUser prompts the user for registration details and registers with the server.
func registerUser(authURL string) bool {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	payload := map[string]string{
		"username": username,
		"password": password,
	}
	data, _ := json.Marshal(payload)

	resp, err := http.Post(fmt.Sprintf("%s/register", authURL), "application/json", strings.NewReader(string(data)))
	if err != nil {
		fmt.Println("Registration failed. Please try again.")
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		fmt.Println("Registration successful!")
		return true
	}
	fmt.Printf("Registration failed with status: %d\n", resp.StatusCode)
	return false
}
