package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/meiti-x/cli_chat/api"
	"github.com/meiti-x/cli_chat/config"
	"github.com/meiti-x/cli_chat/pkg/app_errors"
	c "github.com/meiti-x/cli_chat/pkg/cache"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// FIXME: clear online users in redis on server stop
func main() {
	configPath := flag.String("c", "config.yml", "Path to the configuration file")
	flag.Parse()

	conf, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalln(fmt.Errorf("load config error: %w", err))
	}

	s := api.Setup(conf)
	defer func(Rdb c.Provider) {
		err := Rdb.CloseConnection()
		if err != nil {
			s.Logger.Fatal(c.ErrRedisClose)
		}
	}(s.Rdb)

	http.HandleFunc("/auth/login", api.LoginHandler(s))
	http.HandleFunc("/auth/register", api.RegisterHandler(s))
	http.HandleFunc("/ws", api.InitWS(s))

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", conf.Server.Port),
	}

	go func() {
		fmt.Printf("Server started at %s:%d\n", conf.Server.Host, conf.Server.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.Logger.Error(app_errors.ErrHttpStart, err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer s.Nats.Close()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	fmt.Println("Server gracefully stopped.")
}
