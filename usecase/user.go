package usecase

import (
	"bufio"
	"fmt"
	"github.com/meiti-x/snapp_task/internal/models"
	"github.com/meiti-x/snapp_task/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
)

// AuthenticateUser handles user login and registration.
func AuthenticateUser(userRepo repository.UserRepository) models.Username {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("1. Login")
		fmt.Println("2. Register")
		fmt.Print("Choose an option: ")

		choice, _ := reader.ReadString('\n')
		choice = choice[:len(choice)-1]

		if choice == "1" {
			return loginUser(userRepo)
		} else if choice == "2" {
			return registerUser(userRepo)
		} else {
			fmt.Println("Invalid option. Please try again.")
		}
	}
}

// loginUser handles user login.
func loginUser(userRepo repository.UserRepository) models.Username {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter username: ")
	usernameStr, _ := reader.ReadString('\n')
	var username = models.Username(usernameStr[:len(usernameStr)-1])

	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')
	password = password[:len(password)-1]

	user, err := userRepo.GetUserByUsername(&username)
	if err != nil {
		fmt.Println("User not found. Please register.")
		return AuthenticateUser(userRepo)
	}

	// Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		fmt.Println("Invalid password. Please try again.")
		return AuthenticateUser(userRepo)
	}

	fmt.Println("Login successful!")
	return username
}

// registerUser handles new user registration.
func registerUser(userRepo repository.UserRepository) models.Username {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter username: ")
	usernameStr, _ := reader.ReadString('\n')
	var username = models.Username(usernameStr[:len(usernameStr)-1])

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
	fmt.Println(user)
	if userRepo.CreateUser == nil {
		log.Fatal("UserRepository is nil. Ensure the database is initialized properly.")
	}
	if err := userRepo.CreateUser(&user); err != nil {
		fmt.Printf("Failed to register user: %v\n", err)
		return AuthenticateUser(userRepo)
	}

	fmt.Println("Registration successful!")
	return username
}
