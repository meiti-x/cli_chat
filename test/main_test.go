package test

import (
	"errors"
	"flag"
	"fmt"
	"github.com/meiti-x/snapp_task/api"
	"github.com/meiti-x/snapp_task/config"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	configPath := flag.String("c", "../test.config.yml", "Path to the configuration file")
	flag.Parse()

	conf, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalln(fmt.Errorf("load config error: %w", err))
	}

	s := api.Setup(conf)

	http.HandleFunc("/auth/login", api.LoginHandler(s))
	http.HandleFunc("/auth/register", api.RegisterHandler(s))
	http.HandleFunc("/ws", api.InitWS(s))

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", conf.Server.Port),
	}
	fmt.Println(conf.Server)

	go func() {
		fmt.Printf("Server started at %s:%d\n", conf.Server.Host, conf.Server.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for the server to be ready before running the tests
	waitForServerReady(fmt.Sprintf("http://%s:%d", conf.Server.Host, conf.Server.Port))

	// Run the tests
	exitCode := m.Run()

	// Cleanup
	s.Nats.Close()

	// todo drop db
	// todo clear logger
	server.Shutdown(nil)
	os.Exit(exitCode)
}

func waitForServerReady(url string) {
	// Retry up to 10 times
	for i := 0; i < 10; i++ {
		resp, err := http.Get(url)
		if err == nil {
			resp.Body.Close()
			fmt.Println("Server is ready", url)
			return // Server is ready
		}
		fmt.Printf("Waiting for server to be ready... (error: %v)\n", err)
		time.Sleep(1 * time.Second) // Wait 1 second before retrying
	}
	log.Fatalln("Server failed to start after waiting")
}
